package entities

type tokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type tokenResponse struct {
	Status  bool              `json:"status"`
	Message string            `json:"message"`
	Data    tokenResponseData `json:"data"`
}

type tokenResponseData struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

func NewTokenRequest() tokenRequest {
	return tokenRequest{}
}

func NewTokenResponse(token string, refreshToken string) tokenResponse {
	data := tokenResponseData{Token: token, RefreshToken: refreshToken}
	return tokenResponse{Status: true, Message: operationSuccessMessage, Data: data}
}
