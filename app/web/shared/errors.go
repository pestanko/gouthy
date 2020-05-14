package shared

type ApiError struct {
	Status  string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewApiError(err string, message string, code int) ApiError {
	return ApiError{
		Status:  err,
		Message: message,
		Code:    code,
	}
}

func NotFound(msg string) ApiError {
	return ApiError{
		Status:  "not_found",
		Message: msg,
		Code:    404,
	}
}

func BadRequest(msg string) ApiError {
	return ApiError{
		Status:  "invalid_request",
		Message: msg,
		Code:    400,
	}
}
