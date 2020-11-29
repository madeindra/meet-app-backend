package pubsub

import (
	"encoding/json"
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

func newPubSub() *pubSub {
	return &pubSub{
		Clients:       make([]client, 0),
		Subscriptions: make([]subscription, 0),
	}
}

func (ps *pubSub) addClient(client client) *pubSub {
	ps.Clients = append(ps.Clients, client)
	return ps
}

func (ps *pubSub) removeClient(client client) *pubSub {
	for index, sub := range ps.Subscriptions {
		if client.ID == sub.Client.ID {
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
		}
	}

	for index, c := range ps.Clients {
		if c.ID == client.ID {
			ps.Clients = append(ps.Clients[:index], ps.Clients[index+1:]...)
		}
	}

	return ps
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
		sub.Client.Send(message)
	}
}

func (ps *pubSub) unsubscribe(client *client, topic string) *pubSub {
	for index, sub := range ps.Subscriptions {
		if sub.Client.ID == client.ID && sub.Topic == topic {
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
		}
	}

	return ps
}

func (ps *pubSub) handleReceiveMessage(client client, messageType int, payload []byte) *pubSub {
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
