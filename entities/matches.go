package entities

import "github.com/madeindra/meet-app/models"

type matchRequest struct {
	UserID    uint64 `json:"userId"`
	UserMatch uint64 `json:"userMatch"`
	Liked     bool   `json:"liked"`
}

type matchResponse struct {
	Status  bool              `json:"status"`
	Message string            `json:"message"`
	Data    matchResponseData `json:"data"`
}

type matchBatchResponse struct {
	Status  bool                `json:"status"`
	Message string              `json:"message"`
	Data    []matchResponseData `json:"data"`
}

type matchResponseData struct {
	ID        uint64 `json:"id"`
	UserID    uint64 `json:"userId"`
	UserMatch uint64 `json:"userMatch"`
	Liked     bool   `json:"liked"`
}

func NewMatchRequest() matchRequest {
	return matchRequest{}
}

func NewMatchResponse(id uint64, userID uint64, userMatch uint64, liked bool) matchResponse {
	data := matchResponseData{ID: id, UserID: userID, UserMatch: userMatch, Liked: liked}
	return matchResponse{Status: true, Message: operationSuccessMessage, Data: data}
}

func NewMatchBatchResponse(matches []models.Matches) matchBatchResponse {
	data := []matchResponseData{}

	for i := range matches {
		data = append(data, matchResponseData{ID: matches[i].ID, UserID: matches[i].UserID, UserMatch: matches[i].UserMatch, Liked: matches[i].Liked})
	}

	return matchBatchResponse{Status: true, Message: operationSuccessMessage, Data: data}
}
