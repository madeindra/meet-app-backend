package responses

const (
	PingMessage = "Server is working properly"
)

type pingResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func NewPingResponse() pingResponse {
	return pingResponse{Status: true, Message: PingMessage}
}
