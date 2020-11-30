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

func (ps *pubSub) AddClient(client client) *pubSub {
	ps.Clients = append(ps.Clients, client)
	return ps
}

func (ps *pubSub) RemoveClient(client client) *pubSub {
	for i := 0; i < len(ps.Subscriptions); i++ {
		sub := ps.Subscriptions[i]
		if client.ID == sub.Client.ID {
			if i == len(ps.Subscriptions)-1 {
				ps.Subscriptions = ps.Subscriptions[:len(ps.Subscriptions)-1]
			} else {
				ps.Subscriptions = append(ps.Subscriptions[:i], ps.Subscriptions[i+1:]...)
				i--
			}
		}
	}

	for i := 0; i < len(ps.Clients); i++ {
		c := ps.Clients[i]
		if c.ID == client.ID {
			if i == len(ps.Clients)-1 {
				ps.Clients = ps.Clients[:len(ps.Clients)-1]
			} else {
				ps.Clients = append(ps.Clients[:i], ps.Clients[i+1:]...)
				i--
			}
		}
	}

	return ps
}

func (ps *pubSub) HandleReceiveMessage(client client, messageType int, payload []byte) *pubSub {
	m := newMessage()

	if err := json.Unmarshal(payload, &m); err != nil {
		return ps
	}

	switch m.Action {
	case publish:
		ps.publish(m.Topic, m.Message, nil)
		break

	case subscribe:
		ps.subscribe(&client, m.Topic)
		break

	case unsubscribe:
		ps.unsubscribe(&client, m.Topic)
		break

	default:
		break
	}

	return ps
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
