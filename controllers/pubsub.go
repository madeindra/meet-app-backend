package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/madeindra/meet-app/models"
)

const (
	publish     = "publish"
	subscribe   = "subscribe"
	unsubscribe = "unsubscribe"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func autoID() string {
	return uuid.Must(uuid.NewRandom()).String()
}

type PubSubController struct {
	pubsub models.PubSubInterface
	chat   models.ChatsInterface
}

func NewPubSubController(pubsub models.PubSubInterface, chat models.ChatsInterface) *PubSubController {
	return &PubSubController{pubsub, chat}
}

func (controller *PubSubController) WebsocketHandler(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := controller.pubsub.NewClient(autoID(), conn)
	controller.pubsub.AddClient(client)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			controller.pubsub.RemoveClient(client)
			return
		}
		controller.processMessage(client, messageType, p)
	}
}

func (controller *PubSubController) processMessage(client models.Client, messageType int, payload []byte) error {
	m := controller.pubsub.NewMessage()
	if err := json.Unmarshal(payload, &m); err != nil {
		return errors.New("Failed binding message")
	}

	switch m.Action {
	case publish:
		controller.pubsub.Publish(m.Topic, m.Data, nil)

		ch := controller.chat.New()
		if err := json.Unmarshal(m.Data, &ch); err != nil {
			return errors.New("Failed binding message content")
		}

		if ch.Sender == 0 || ch.Target == 0 || ch.Content == "" {
			log.Printf("Bad message %s", string(m.Data))
			return errors.New("Message is not in a proper format")
		}

		go controller.chat.Create(ch)
		break

	case subscribe:
		controller.pubsub.Subscribe(&client, m.Topic)
		break

	case unsubscribe:
		controller.pubsub.Unsubscribe(&client, m.Topic)
		break

	default:
		break
	}

	return nil
}
