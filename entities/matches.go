package entities

import "github.com/madeindra/meet-app/models"

type matchRequest struct {
	UserID    uint64 `json:"userId"`
	UserMatch uint64 `json:"userMatch"`
	Matched   bool   `json:"matched"`
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
	Matched   bool   `json:"matched"`
}

func NewMatchRequest() matchRequest {
	return matchRequest{}
}

func NewMatchResponse(id uint64, userID uint64, userMatch uint64, matched bool) matchResponse {
	data := matchResponseData{ID: id, UserID: userID, UserMatch: userMatch, Matched: matched}
	return matchResponse{Status: true, Message: operationSuccessMessage, Data: data}
}

func NewMatchBatchResponse(matches []models.Matches) matchBatchResponse {
	data := []matchResponseData{}

	for i := range matches {
		data = append(data, matchResponseData{ID: matches[i].ID, UserID: matches[i].UserID, UserMatch: matches[i].UserMatch, Matched: matches[i].Matched})
	}

	return matchBatchResponse{Status: true, Message: operationSuccessMessage, Data: data}
}
