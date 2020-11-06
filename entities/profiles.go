package entities

import "github.com/madeindra/meet-app/models"

type profileRequest struct {
	UserID      uint64  `json:"userId"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type profileResponse struct {
	Status  bool                `json:"status"`
	Message string              `json:"message"`
	Data    profileResponseData `json:"data"`
}

type profileBatchResponse struct {
	Status  bool                  `json:"status"`
	Message string                `json:"message"`
	Data    []profileResponseData `json:"data"`
}

type profileResponseData struct {
	UserID      uint64  `json:"userId"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

func NewProfileRequest() profileRequest {
	return profileRequest{}
}

func NewProfileResponse(userId uint64, firstName string, lastName string, description string, latitude float64, longitude float64) profileResponse {
	data := profileResponseData{UserID: userId, FirstName: firstName, LastName: lastName, Description: description, Latitude: latitude, Longitude: longitude}
	return profileResponse{Status: true, Message: operationSuccessMessage, Data: data}
}

func NewProfileBatchResponse(profiles []models.Profiles) profileBatchResponse {
	data := []profileResponseData{}
	for i := range profiles {
		data = append(data, profileResponseData{UserID: profiles[i].UserID, FirstName: profiles[i].FirstName, LastName: profiles[i].LastName, Description: profiles[i].Description, Latitude: profiles[i].Latitude, Longitude: profiles[i].Longitude})
	}
	return profileBatchResponse{Status: true, Message: operationSuccessMessage, Data: data}
}
