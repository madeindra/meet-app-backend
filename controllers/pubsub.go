package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/madeindra/meet-app/entities"
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

type PubSubController struct {
	pubsub     models.PubSubInterface
	chat       models.ChatsInterface
	ticket     models.TicketInterface
	credential models.CredentialInterface
}

func NewPubSubController(pubsub models.PubSubInterface, chat models.ChatsInterface, ticket models.TicketInterface, credential models.CredentialInterface) *PubSubController {
	return &PubSubController{pubsub, chat, ticket, credential}
}

func (controller *PubSubController) WebsocketHandler(ctx *gin.Context) {
	userID, err := controller.getTicketUser(ctx.Query("ticket"))
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := controller.pubsub.NewClient(userID, conn)
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

func (controller *PubSubController) processMessage(client models.Client, messageType int, payload []byte) {
	m := controller.pubsub.NewMessage()
	if err := json.Unmarshal(payload, &m); err != nil {
		controller.pubsub.BounceBack(&client, "Server: Failed binding message")
	}

	switch m.Action {
	case publish:
		ch := controller.chat.New()
		if err := json.Unmarshal(m.Data, &ch); err != nil {
			controller.pubsub.BounceBack(&client, "Server: Failed binding message content")
			break
		}

		ch.SenderID = client.ID
		if ch.TargetID == 0 || ch.Content == "" {
			controller.pubsub.BounceBack(&client, "Server: Message is not in a proper format")
			break
		}

		user := controller.credential.New()
		user.ID = ch.TargetID
		if existing := controller.credential.FindOne(user); existing.Email == "" {
			controller.pubsub.BounceBack(&client, "Server: Target does not exist")
			break
		}

		chat := entities.NewChatResponse(ch.ID, ch.SenderID, ch.TargetID, ch.Content)

		res, err := json.Marshal(&chat)
		if err != nil {
			controller.pubsub.BounceBack(&client, "Server: Failed creating raw message")
			break
		}

		go controller.pubsub.Publish(ch.TargetID, res, nil)
		go controller.chat.Create(ch)
		break

	case subscribe:
		controller.pubsub.Subscribe(&client, client.ID)
		break

	case unsubscribe:
		controller.pubsub.Unsubscribe(&client, client.ID)
		break

	default:
		controller.pubsub.BounceBack(&client, "Server: Action unrecognized")
		break
	}
}

func (controller *PubSubController) getTicketUser(ticket string) (uint64, error) {
	if ticket == "" {
		return 0, errors.New("Invalid ticket")
	}

	data := controller.ticket.New()
	data.Ticket = ticket

	existing := controller.ticket.FindOne(data)
	if existing.ID == 0 {
		return 0, errors.New("Ticket not found")
	}

	return existing.UserID, nil
}
