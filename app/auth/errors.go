package auth

import (
	"github.com/pestanko/gouthy/app/shared"
)

func NewOAuth2AuthorizationRequestError(msg string) shared.AppError {
	return shared.NewAppError(msg).WithType("oauth2_auth_request")
}

