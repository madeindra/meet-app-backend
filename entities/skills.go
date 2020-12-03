package entities

import "github.com/madeindra/meet-app/models"

type skillRequest struct {
	UserID    uint64 `json:"userId"`
	SkillName string `json:"skillName"`
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
	ID        uint64 `json:"id"`
	UserID    uint64 `json:"userId"`
	SkillName string `json:"skillName"`
}

func NewSkillRequest() skillRequest {
	return skillRequest{}
}

func NewSkillResponse(id uint64, userId uint64, skillName string) skillResponse {
	data := skillResponseData{ID: id, UserID: userId, SkillName: skillName}
	return skillResponse{Status: true, Message: operationSuccessMessage, Data: data}
}

func NewSkillBatchResponse(skills []models.Skills) skillBatchResponse {
	data := []skillResponseData{}
	for i := range skills {
		data = append(data, skillResponseData{ID: skills[i].ID, UserID: skills[i].UserID, SkillName: skills[i].SkillName})
	}
	return skillBatchResponse{Status: true, Message: operationSuccessMessage, Data: data}
}
