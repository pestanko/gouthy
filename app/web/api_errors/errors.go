package api_errors

import "github.com/gin-gonic/gin"

type ApiError gin.H
type ErrorDetail gin.H

func NewApiError() ApiError {
	return ApiError{
		"code":    500,
		"status":  "server_error",
	}
}

func (e ApiError) Code() int {
	return e["code"].(int)
}

func (e ApiError) WithCode(code int) ApiError {
	e["code"] = code
	return e
}

func (e ApiError) WithMessage(message string) ApiError {
	e["message"] = message
	return e
}

func (e ApiError) WithStatus(status string) ApiError {
	e["status"] = status
	return e
}

func (e ApiError) WithError(err error) ApiError {
	e["error"] = err.Error()
	return e
}

func (e ApiError) WithDetail(detail ErrorDetail) ApiError {
	e["detail"] = detail
	return e
}

func NewNotFound() ApiError {
	return NewApiError().WithStatus("not_found").WithCode(404)
}

func NewInvalidRequest() ApiError {
	return NewApiError().WithStatus("invalid_request").WithCode(400)
}
func NewUnauthorizedError() ApiError {
	return NewApiError().WithStatus("unauthorized").WithCode(401)
}
