package pubsub

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func autoID() string {
	return uuid.Must(uuid.NewRandom()).String()
}

var ps = newPubSub()

func WebsocketHandler(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(autoID(), conn)
	ps.addClient(client)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			ps.removeClient(client)
			return
		}
		ps.handleReceiveMessage(client, messageType, p)
	}
}
