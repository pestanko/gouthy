package auth

import (
	"github.com/pestanko/gouthy/app/shared"
)

func NewOAuth2AuthorizationRequestError(msg string) shared.GouthyError {
	return shared.NewGouthyError(msg).WithType("oauth2_auth_request")
}

