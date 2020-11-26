package entities

type resetStartRequest struct {
	Email string `json:"email" binding:"required"`
}

type resetCompleteRequest struct {
	Password string `json:"password" binding:"required"`
}

type resetStartResponse struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    resetStartData `json:"data"`
}

type resetCompleteResponse struct {
	Status  bool              `json:"status"`
	Message string            `json:"message"`
	Data    resetCompleteData `json:"data"`
}

type resetStartData struct {
	UserID uint64 `json:"userId"`
	Token  string `json:"token"`
}

type resetCompleteData struct {
	UserID uint64 `json:"userId"`
}

func NewResetStartRequest() resetStartRequest {
	return resetStartRequest{}
}

func NewResetCompleteRequest() resetCompleteRequest {
	return resetCompleteRequest{}
}

func NewResetStartResponse(userId uint64, token string) resetStartResponse {
	data := resetStartData{UserID: userId, Token: token}
	return resetStartResponse{Status: true, Message: operationSuccessMessage, Data: data}
}

func NewResetCompleteResponse(userId uint64) resetCompleteResponse {
	data := resetCompleteData{UserID: userId}
	return resetCompleteResponse{Status: true, Message: operationSuccessMessage, Data: data}
}
