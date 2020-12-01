package models

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
)

const (
	publish     = "publish"
	subscribe   = "subscribe"
	unsubscribe = "unsubscribe"
)

type pubSub struct {
	Clients       []client
	Subscriptions []subscription
}

type PubSubInterface interface {
	AddClient(client client) *pubSub
	AddChat(data Chats) (Chats, error)
	RemoveClient(client client) *pubSub
	HandleReceiveMessage(client client, messageType int, payload []byte) (*pubSub, error)
}

type PubSubImplementation struct {
	db     *gorm.DB
	pubSub *pubSub
}

type client struct {
	ID         string
	Connection *websocket.Conn
}

type subscription struct {
	Topic  string
	Client *client
}

type message struct {
	Action string          `json:"action"`
	Topic  string          `json:"topic"`
	Data   json.RawMessage `json:"data"`
}

func NewPubSub() *pubSub {
	return &pubSub{
		Clients:       make([]client, 0),
		Subscriptions: make([]subscription, 0),
	}
}

func NewPubSubModel(db *gorm.DB, ps *pubSub) *PubSubImplementation {
	return &PubSubImplementation{db: db, pubSub: ps}
}

func NewClient(id string, conn *websocket.Conn) client {
	return client{
		ID:         id,
		Connection: conn,
	}
}

func newMessageContent() Chats {
	return Chats{}
}

func newSubscription(topic string, client *client) subscription {
	return subscription{
		Topic:  topic,
		Client: client,
	}
}

func newMessage() *message {
	return &message{}
}

func (implementation *PubSubImplementation) AddClient(client client) *pubSub {
	implementation.pubSub.Clients = append(implementation.pubSub.Clients, client)
	return implementation.pubSub
}

func (implementation *PubSubImplementation) RemoveClient(client client) *pubSub {
	for i := 0; i < len(implementation.pubSub.Subscriptions); i++ {
		sub := implementation.pubSub.Subscriptions[i]
		if client.ID == sub.Client.ID {
			if i == len(implementation.pubSub.Subscriptions)-1 {
				implementation.pubSub.Subscriptions = implementation.pubSub.Subscriptions[:len(implementation.pubSub.Subscriptions)-1]
			} else {
				implementation.pubSub.Subscriptions = append(implementation.pubSub.Subscriptions[:i], implementation.pubSub.Subscriptions[i+1:]...)
				i--
			}
		}
	}

	for i := 0; i < len(implementation.pubSub.Clients); i++ {
		c := implementation.pubSub.Clients[i]
		if c.ID == client.ID {
			if i == len(implementation.pubSub.Clients)-1 {
				implementation.pubSub.Clients = implementation.pubSub.Clients[:len(implementation.pubSub.Clients)-1]
			} else {
				implementation.pubSub.Clients = append(implementation.pubSub.Clients[:i], implementation.pubSub.Clients[i+1:]...)
				i--
			}
		}
	}

	return implementation.pubSub
}

func (implementation *PubSubImplementation) HandleReceiveMessage(client client, messageType int, payload []byte) (*pubSub, error) {
	m := newMessage()
	if err := json.Unmarshal(payload, &m); err != nil {
		return implementation.pubSub, errors.New("Failed binding message")
	}

	switch m.Action {
	case publish:
		implementation.pubSub.publish(m.Topic, m.Data, nil)

		ch := newMessageContent()
		if err := json.Unmarshal(m.Data, &ch); err != nil {
			return implementation.pubSub, errors.New("Failed binding message content")
		}

		if ch.Sender == 0 || ch.Target == 0 || ch.Content == "" {
			log.Printf("Bad message %s", string(m.Data))
			return implementation.pubSub, errors.New("Message is not in a proper format")
		}

		go implementation.AddChat(ch)
		break

	case subscribe:
		implementation.pubSub.subscribe(&client, m.Topic)
		break

	case unsubscribe:
		implementation.pubSub.unsubscribe(&client, m.Topic)
		break

	default:
		break
	}

	return implementation.pubSub, nil
}

func (implementation *PubSubImplementation) AddChat(data Chats) (Chats, error) {
	tx := implementation.db.Begin()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return Chats{}, err
	}

	return data, tx.Commit().Error
}

func (client *client) send(message []byte) error {
	return client.Connection.WriteMessage(1, message)
}

func (ps *pubSub) getSubscriptions(topic string, client *client) []subscription {
	var subscriptionList []subscription

	for _, subscription := range ps.Subscriptions {
		if client != nil {
			if subscription.Client.ID == client.ID && subscription.Topic == topic {
				subscriptionList = append(subscriptionList, subscription)
			}

		} else {
			if subscription.Topic == topic {
				subscriptionList = append(subscriptionList, subscription)
			}
		}
	}

	return subscriptionList
}

func (ps *pubSub) subscribe(client *client, topic string) *pubSub {
	clientSubs := ps.getSubscriptions(topic, client)
	if len(clientSubs) > 0 {
		return ps
	}

	subsctiption := newSubscription(topic, client)
	ps.Subscriptions = append(ps.Subscriptions, subsctiption)
	return ps
}

func (ps *pubSub) publish(topic string, message []byte, excludeClient *client) {
	subscriptions := ps.getSubscriptions(topic, nil)

	for _, sub := range subscriptions {
		sub.Client.send(message)
	}
}

func (ps *pubSub) unsubscribe(client *client, topic string) *pubSub {
	for i := 0; i < len(ps.Subscriptions); i++ {
		sub := ps.Subscriptions[i]
		if sub.Client.ID == client.ID && sub.Topic == topic {
			if i == len(ps.Subscriptions)-1 {
				ps.Subscriptions = ps.Subscriptions[:len(ps.Subscriptions)-1]
			} else {
				ps.Subscriptions = append(ps.Subscriptions[:i], ps.Subscriptions[i+1:]...)
				i--
			}
		}
	}

	return ps
}
