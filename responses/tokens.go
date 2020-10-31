package responses

type tokenResponse struct {
	Status  bool      `json:"status"`
	Message string    `json:"message"`
	Data    tokenData `json:"data"`
}

type tokenData struct {
	UserID       uint64 `json:"userId"`
	RefreshToken string `json:"refreshToken"`
}

func NewTokenResponse(userID uint64, refreshToken string) tokenResponse {
	data := tokenData{UserID: userID, RefreshToken: refreshToken}
	return tokenResponse{Status: true, Message: operationSuccessMessage, Data: data}
}
