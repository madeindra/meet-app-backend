package entities

type tokenResponse struct {
	Status  bool      `json:"status"`
	Message string    `json:"message"`
	Data    tokenData `json:"data"`
}

type tokenData struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

func NewTokenResponse(token string, refreshToken string) tokenResponse {
	data := tokenData{Token: token, RefreshToken: refreshToken}
	return tokenResponse{Status: true, Message: operationSuccessMessage, Data: data}
}
