package web_utils

type ApiError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewApiError(err string, message string, code int) ApiError {
	return ApiError{
		Error:   err,
		Message: message,
		Code:    code,
	}
}