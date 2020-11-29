package pubsub

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

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

func newClient(id string, conn *websocket.Conn) client {
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

func (client *client) Send(message []byte) error {
	return client.Connection.WriteMessage(1, message)
}
