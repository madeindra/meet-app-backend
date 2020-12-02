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
	NewClient(id uint64, conn *websocket.Conn) Client
	NewMessage() *message
	AddClient(client Client) *pubSub
	RemoveClient(client Client) *pubSub
	Publish(topic uint64, message []byte, excludeClient *Client)
	BounceBack(client *Client, message string)
	Subscribe(client *Client, topic uint64) *pubSub
	Unsubscribe(client *Client, topic uint64) *pubSub
}

type PubSubImplementation struct {
	pubSub *pubSub
}

type Client struct {
	ID         uint64
	Connection *websocket.Conn
}

type subscription struct {
	Topic  uint64
	Client *Client
}

type message struct {
	Action string          `json:"action"`
	Topic  uint64          `json:"topic"`
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

func (implementation *PubSubImplementation) NewClient(id uint64, conn *websocket.Conn) Client {
	return Client{
		ID:         id,
		Connection: conn,
	}
}

func (implementation *PubSubImplementation) NewMessage() *message {
	return &message{}
}

func (implementation *PubSubImplementation) AddClient(client Client) *pubSub {
	implementation.pubSub.Clients = append(implementation.pubSub.Clients, client)
	return implementation.pubSub
}

func (implementation *PubSubImplementation) RemoveClient(client Client) *pubSub {
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

func (implementation *PubSubImplementation) Publish(topic uint64, message []byte, excludeClient *Client) {
	subscriptions := implementation.pubSub.getSubscriptions(topic, nil)

	for _, sub := range subscriptions {
		sub.Client.send(message)
	}
}

func (implementation *PubSubImplementation) BounceBack(client *Client, message string) {
	client.send([]byte(message))
}

func (implementation *PubSubImplementation) Subscribe(client *Client, topic uint64) *pubSub {
	clientSubs := implementation.pubSub.getSubscriptions(topic, client)
	if len(clientSubs) > 0 {
		return implementation.pubSub
	}

	subscription := implementation.pubSub.newSubscription(topic, client)
	implementation.pubSub.Subscriptions = append(implementation.pubSub.Subscriptions, subscription)
	return implementation.pubSub
}

func (implementation *PubSubImplementation) Unsubscribe(client *Client, topic uint64) *pubSub {
	for i := 0; i < len(implementation.pubSub.Subscriptions); i++ {
		sub := implementation.pubSub.Subscriptions[i]
		if sub.Client.ID == client.ID && sub.Topic == topic {
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

func (ps *pubSub) newSubscription(topic uint64, Client *Client) subscription {
	return subscription{
		Topic:  topic,
		Client: Client,
	}
}

func (ps *pubSub) getSubscriptions(topic uint64, Client *Client) []subscription {
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
