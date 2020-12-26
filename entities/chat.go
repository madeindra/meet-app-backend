package entities

import "github.com/madeindra/meet-app/models"

type chatBatchResponse struct {
	Status  bool               `json:"status"`
	Message string             `json:"message"`
	Data    []chatResponseData `json:"data"`
}

type chatResponseData struct {
	ID      uint64 `gorm:"primaryKey" json:"id,omitempty"`
	Sender  uint64 `json:"senderId"`
	Target  uint64 `json:"targetId"`
	Content string `json:"content"`
}

func NewChatBatchResponse(chats []models.Chats) chatBatchResponse {
	data := []chatResponseData{}
	for i := range chats {
		data = append(data, chatResponseData{ID: chats[i].ID, Sender: chats[i].Sender, Target: chats[i].Target, Content: chats[i].Content})
	}
	return chatBatchResponse{Status: true, Message: operationSuccessMessage, Data: data}
}
