package entities

const (
	registerSuccessMessage     = "Registration Successful"
	authenticateSuccessMessage = "Authentication Successful"
)

type credentialRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type credentialUpdateRequest struct {
	Email       string `json:"email"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword" binding:"required_with=OldPassword"`
}

type credentialResponse struct {
	Status  bool                   `json:"status"`
	Message string                 `json:"message"`
	Data    credentialResponseData `json:"data"`
}

type authenticatedResponse struct {
	Status  bool                      `json:"status"`
	Message string                    `json:"message"`
	Data    authenticatedResponseData `json:"data"`
}

type credentialResponseData struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
}

type authenticatedResponseData struct {
	ID           uint64 `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

func NewCredentialRequest() credentialRequest {
	return credentialRequest{}
}

func NewCredentialUpdateRequest() credentialUpdateRequest {
	return credentialUpdateRequest{}
}

func NewCredentialResponse(ID uint64, email string) credentialResponse {
	data := credentialResponseData{ID: ID, Email: email}
	return credentialResponse{Status: true, Message: registerSuccessMessage, Data: data}
}

func NewCredentialUpdateResponse(ID uint64, email string) credentialResponse {
	data := credentialResponseData{ID: ID, Email: email}
	return credentialResponse{Status: true, Message: operationSuccessMessage, Data: data}
}

func NewAuthenticatedResponse(ID uint64, email string, token string, refreshToken string) authenticatedResponse {
	data := authenticatedResponseData{ID: ID, Email: email, Token: token, RefreshToken: refreshToken}
	return authenticatedResponse{Status: true, Message: authenticateSuccessMessage, Data: data}
}
