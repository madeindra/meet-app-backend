package responses

const (
	badRequestMessage           string = "Bad Request"
	interenalServerErrorMessage string = "Internal Server Error"
	notFoundMessage             string = "Not Found"
	unauthorizedMessage         string = "Unauthorized"
)

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

func NotFoundResponse() failed {
	return failedResponse(notFoundMessage)
}

func UnauthorizedResponse() failed {
	return failedResponse(unauthorizedMessage)
}

func InterenalServerErrorResponse() failed {
	return failedResponse(interenalServerErrorMessage)
}
