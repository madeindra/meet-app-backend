package responses

const (
	registerSuccessMessage = "Registration Successful"
)

type credentialResponse struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    credentialData `json:"data"`
}

type credentialData struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

func NewCredentialResponse(ID uint, email string) credentialResponse {
	data := credentialData{ID: ID, Email: email}
	return credentialResponse{Status: true, Message: registerSuccessMessage, Data: data}
}
