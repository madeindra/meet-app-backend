package responses

type credentialResponse struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    credentialData `json:"data"`
}

type credentialData struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

func NewCredentialData(ID uint, email string) credentialData {
	return credentialData{ID: ID, Email: email}
}

func NewCredentialResponse(data credentialData) credentialResponse {
	return credentialResponse{Status: true, Message: "Registration Successful", Data: data}
}
