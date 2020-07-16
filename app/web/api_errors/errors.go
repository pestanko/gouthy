package api_errors

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/shared"
)

type ApiError gin.H
type ErrorDetail map[string]interface{}

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

func (e ApiError) WithEntity(entity string) ApiError {
	e["entity"] = entity
	return e
}

func FromGouthyError(err shared.GouthyError) ApiError {
	return NewApiError().WithDetail(map[string]interface{}(err.Detail()))
}

/**
 Not Found
 */

func NewNotFound() ApiError {
	return NewApiError().WithStatus("not_found").WithCode(404)
}

func NewUserNotFound() ApiError {
	return NewNotFound().WithMessage("User not found").WithEntity("user")
}

func NewAppNotFound() ApiError {
	return NewNotFound().WithMessage("Application not found").WithEntity("application")
}

/**
 * Bad request
 */

func NewInvalidRequest() ApiError {
	return NewApiError().WithStatus("invalid_request").WithCode(400)
}
func NewUnauthorizedError() ApiError {
	return NewApiError().WithStatus("unauthorized").WithCode(401)
}
