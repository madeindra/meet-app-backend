package entities

import "github.com/madeindra/meet-app/models"

type skillRequest struct {
	UserID uint64 `json:"userId"`
	Name   string `json:"name"`
}

type skillResponse struct {
	Status  bool              `json:"status"`
	Message string            `json:"message"`
	Data    skillResponseData `json:"data"`
}

type skillBatchResponse struct {
	Status  bool                `json:"status"`
	Message string              `json:"message"`
	Data    []skillResponseData `json:"data"`
}

type skillResponseData struct {
	UserID uint64 `json:"userId"`
	Name   string `json:"name"`
}

func NewSkillRequest() skillRequest {
	return skillRequest{}
}

func NewSkillResponse(userId uint64, name string) skillResponse {
	data := skillResponseData{UserID: userId, Name: name}
	return skillResponse{Status: true, Message: operationSuccessMessage, Data: data}
}

func NewSkillBatchResponse(skills []models.Skills) skillBatchResponse {
	data := []skillResponseData{}
	for i := range skills {
		data = append(data, skillResponseData{UserID: skills[i].UserID, Name: skills[i].Name})
	}
	return skillBatchResponse{Status: true, Message: operationSuccessMessage, Data: data}
}
