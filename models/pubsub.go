package models

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type pubSub struct {
	Clients       []Client
	Subscriptions []subscription
}

type PubSubInterface interface {
	NewClient(id string, conn *websocket.Conn) Client
	NewMessage() *message
	AddClient(Client Client) *pubSub
	RemoveClient(Client Client) *pubSub
	Publish(topic string, message []byte, excludeClient *Client)
	Subscribe(Client *Client, topic string) *pubSub
	Unsubscribe(Client *Client, topic string) *pubSub
}

type PubSubImplementation struct {
	pubSub *pubSub
}

type Client struct {
	ID         string
	Connection *websocket.Conn
}

type subscription struct {
	Topic  string
	Client *Client
}

type message struct {
	Action string          `json:"action"`
	Topic  string          `json:"topic"`
	Data   json.RawMessage `json:"data"`
}

func NewPubSub() *pubSub {
	return &pubSub{
		Clients:       make([]Client, 0),
		Subscriptions: make([]subscription, 0),
	}
}

func NewPubSubModel(ps *pubSub) *PubSubImplementation {
	return &PubSubImplementation{pubSub: ps}
}

func (implementation *PubSubImplementation) NewClient(id string, conn *websocket.Conn) Client {
	return Client{
		ID:         id,
		Connection: conn,
	}
}

func (implementation *PubSubImplementation) NewMessage() *message {
	return &message{}
}

func (implementation *PubSubImplementation) AddClient(Client Client) *pubSub {
	implementation.pubSub.Clients = append(implementation.pubSub.Clients, Client)
	return implementation.pubSub
}

func (implementation *PubSubImplementation) RemoveClient(Client Client) *pubSub {
	for i := 0; i < len(implementation.pubSub.Subscriptions); i++ {
		sub := implementation.pubSub.Subscriptions[i]
		if Client.ID == sub.Client.ID {
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
		if c.ID == Client.ID {
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

func (implementation *PubSubImplementation) Publish(topic string, message []byte, excludeClient *Client) {
	subscriptions := implementation.pubSub.getSubscriptions(topic, nil)

	for _, sub := range subscriptions {
		sub.Client.send(message)
	}
}

func (implementation *PubSubImplementation) Subscribe(Client *Client, topic string) *pubSub {
	clientSubs := implementation.pubSub.getSubscriptions(topic, Client)
	if len(clientSubs) > 0 {
		return implementation.pubSub
	}

	subscription := implementation.pubSub.newSubscription(topic, Client)
	implementation.pubSub.Subscriptions = append(implementation.pubSub.Subscriptions, subscription)
	return implementation.pubSub
}

func (implementation *PubSubImplementation) Unsubscribe(Client *Client, topic string) *pubSub {
	for i := 0; i < len(implementation.pubSub.Subscriptions); i++ {
		sub := implementation.pubSub.Subscriptions[i]
		if sub.Client.ID == Client.ID && sub.Topic == topic {
			if i == len(implementation.pubSub.Subscriptions)-1 {
				implementation.pubSub.Subscriptions = implementation.pubSub.Subscriptions[:len(implementation.pubSub.Subscriptions)-1]
			} else {
				implementation.pubSub.Subscriptions = append(implementation.pubSub.Subscriptions[:i], implementation.pubSub.Subscriptions[i+1:]...)
				i--
			}
		}
	}

	return implementation.pubSub
}

func (Client *Client) send(message []byte) error {
	return Client.Connection.WriteMessage(1, message)
}

func (ps *pubSub) newSubscription(topic string, Client *Client) subscription {
	return subscription{
		Topic:  topic,
		Client: Client,
	}
}

func (ps *pubSub) getSubscriptions(topic string, Client *Client) []subscription {
	var subscriptionList []subscription

	for _, subscription := range ps.Subscriptions {
		if Client != nil {
			if subscription.Client.ID == Client.ID && subscription.Topic == topic {
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
