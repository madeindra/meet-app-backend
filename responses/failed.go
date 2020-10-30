package responses

const (
	badRequestMessage           = "Bad Request"
	interenalServerErrorMessage = "Internal Server Error"
	notFoundMessage             = "Not Found"
	unauthorizedMessage         = "Unauthorized"
	conflictMessage             = "Conflict"
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

func ConflictResponse() failed {
	return failedResponse(conflictMessage)
}
