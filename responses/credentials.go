package responses

const (
	registerSuccessMessage     = "Registration Successful"
	authenticateSuccessMessage = "Authentication Successful"
)

type credentialResponse struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    credentialData `json:"data"`
}

type authenticatedResponse struct {
	Status  bool              `json:"status"`
	Message string            `json:"message"`
	Data    authenticatedData `json:"data"`
}

type credentialData struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type authenticatedData struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

func NewCredentialResponse(ID uint, email string) credentialResponse {
	data := credentialData{ID: ID, Email: email}
	return credentialResponse{Status: true, Message: registerSuccessMessage, Data: data}
}

func NewAuthenticatedResponse(ID uint, email string, token string, refreshToken string) authenticatedResponse {
	data := authenticatedData{ID: ID, Email: email, Token: token, RefreshToken: refreshToken}
	return authenticatedResponse{Status: true, Message: authenticateSuccessMessage, Data: data}
}
