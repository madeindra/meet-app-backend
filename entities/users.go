package entities

import "github.com/madeindra/meet-app/models"

type userResponse struct {
	Status  bool             `json:"status"`
	Message string           `json:"message"`
	Data    userResponseData `json:"data"`
}

type userBatchResponse struct {
	Status  bool               `json:"status"`
	Message string             `json:"message"`
	Data    []userResponseData `json:"data"`
}

type userResponseData struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUserResponse(id uint64, name string, email string) userResponse {
	data := userResponseData{ID: id, Name: name, Email: email}
	return userResponse{Status: true, Message: operationSuccessMessage, Data: data}
}

func NewUserBatchResponse(users []models.Profiles) userBatchResponse {
	data := []userResponseData{}
	for i := range users {
		data = append(data, userResponseData{ID: users[i].CredentialID, Name: users[i].Name, Email: users[i].Credential.Email})
	}
	return userBatchResponse{Status: true, Message: operationSuccessMessage, Data: data}
}
