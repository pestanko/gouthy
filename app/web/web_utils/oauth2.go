package web_utils

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/domain/auth"
	"strings"
)

func OAuth2ParseAuthorizationRequest(request *gin.Context) auth.OAuth2AuthorizationRequest {
	scopes := strings.Split(request.Param("scopes"), " ")
	return auth.OAuth2AuthorizationRequest{
		ClientId:      request.Param("client_id"),
		ResponseType:  request.Param("response_type"),
		RedirectUri:   request.Param("request_uri"),
		Scopes:        scopes,
		State:         request.Param("state"),
		Nonce:         request.Param("nonce"),
		PKCEChallenge: request.Param("code_challenge"),
		PKCEMethod:    request.Param("code_challenge_method"),
	}
}
