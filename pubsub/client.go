package pubsub

import "github.com/gorilla/websocket"

type client struct {
	ID         string
	Connection *websocket.Conn
}

type subscription struct {
	Topic  string
	Client *client
}

func (client *client) Send(message []byte) error {
	return client.Connection.WriteMessage(1, message)
}
