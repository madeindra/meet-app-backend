package responses

const badRequestMessage string = "Bad Request"
const interenalServerErrorMessage string = "Internal Server Error"

type failed struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func failedResponse(message string) failed {
	return failed{Status: false, Message: message}
}

func BadRequestResponse() failed {
	return failedResponse(badRequestMessage)
}

func InterenalServerErrorResponse() failed {
	return failedResponse(interenalServerErrorMessage)
}
