package models

import (
	"encoding/json"

	"github.com/gorilla/websocket"
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
	RemoveClient(client client) *pubSub
	HandleReceiveMessage(client client, messageType int, payload []byte) *pubSub
}

type PubSubImplementation struct {
	PubSub *pubSub
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
	Action  string          `json:"action"`
	Topic   string          `json:"topic"`
	Message json.RawMessage `json:"message"`
}

func NewPubSub() *pubSub {
	return &pubSub{
		Clients:       make([]client, 0),
		Subscriptions: make([]subscription, 0),
	}
}

func NewPubSubModel(ps *pubSub) *PubSubImplementation {
	return &PubSubImplementation{PubSub: ps}
}

func NewClient(id string, conn *websocket.Conn) client {
	return client{
		ID:         id,
		Connection: conn,
	}
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
	implementation.PubSub.Clients = append(implementation.PubSub.Clients, client)
	return implementation.PubSub
}

func (implementation *PubSubImplementation) RemoveClient(client client) *pubSub {
	for i := 0; i < len(implementation.PubSub.Subscriptions); i++ {
		sub := implementation.PubSub.Subscriptions[i]
		if client.ID == sub.Client.ID {
			if i == len(implementation.PubSub.Subscriptions)-1 {
				implementation.PubSub.Subscriptions = implementation.PubSub.Subscriptions[:len(implementation.PubSub.Subscriptions)-1]
			} else {
				implementation.PubSub.Subscriptions = append(implementation.PubSub.Subscriptions[:i], implementation.PubSub.Subscriptions[i+1:]...)
				i--
			}
		}
	}

	for i := 0; i < len(implementation.PubSub.Clients); i++ {
		c := implementation.PubSub.Clients[i]
		if c.ID == client.ID {
			if i == len(implementation.PubSub.Clients)-1 {
				implementation.PubSub.Clients = implementation.PubSub.Clients[:len(implementation.PubSub.Clients)-1]
			} else {
				implementation.PubSub.Clients = append(implementation.PubSub.Clients[:i], implementation.PubSub.Clients[i+1:]...)
				i--
			}
		}
	}

	return implementation.PubSub
}

func (implementation *PubSubImplementation) HandleReceiveMessage(client client, messageType int, payload []byte) *pubSub {
	m := newMessage()

	if err := json.Unmarshal(payload, &m); err != nil {
		return implementation.PubSub
	}

	switch m.Action {
	case publish:
		implementation.PubSub.publish(m.Topic, m.Message, nil)
		break

	case subscribe:
		implementation.PubSub.subscribe(&client, m.Topic)
		break

	case unsubscribe:
		implementation.PubSub.unsubscribe(&client, m.Topic)
		break

	default:
		break
	}

	return implementation.PubSub
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
