package entities

import "github.com/madeindra/meet-app/models"

type chatResponse struct {
	Status bool             `json:"status"`
	Data   chatResponseData `json:"data"`
}
type chatBatchResponse struct {
	Status  bool               `json:"status"`
	Message string             `json:"message"`
	Data    []chatResponseData `json:"data"`
}

type chatResponseData struct {
	ID            uint64 `gorm:"primaryKey" json:"id,omitempty"`
	SenderID      uint64 `json:"senderId"`
	SenderName    string `json:"senderName"`
	SenderPicture string `json:"senderPicture"`
	TargetID      uint64 `json:"targetId"`
	Content       string `json:"content"`
}

func NewChatResponse(id uint64, senderID uint64, targetID uint64, content string) chatResponse {
	data := chatResponseData{ID: id, SenderID: senderID, TargetID: targetID, Content: content}
	return chatResponse{Status: true, Data: data}
}

func NewChatBatchResponse(chats []models.Chats) chatBatchResponse {
	data := []chatResponseData{}
	for i := range chats {
		data = append(data, chatResponseData{ID: chats[i].ID, SenderID: chats[i].SenderID, SenderName: chats[i].SenderName, SenderPicture: chats[i].SenderPicture, TargetID: chats[i].TargetID, Content: chats[i].Content})
	}
	return chatBatchResponse{Status: true, Message: operationSuccessMessage, Data: data}
}
