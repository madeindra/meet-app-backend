package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/madeindra/meet-app/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func autoID() string {
	return uuid.Must(uuid.NewRandom()).String()
}

var ps = models.NewPubSub()

func WebsocketHandler(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := models.NewClient(autoID(), conn)
	ps.AddClient(client)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			ps.RemoveClient(client)
			return
		}
		ps.HandleReceiveMessage(client, messageType, p)
	}
}
