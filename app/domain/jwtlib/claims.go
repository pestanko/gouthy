package jwtlib

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/pestanko/gouthy/app/domain/applications"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/shared"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

type Claims struct {
	Audience   string `json:"aud,omitempty"`
	ExpiresAt  int64  `json:"exp,omitempty"`
	Id         string `json:"jti,omitempty"`
	IssuedAt   int64  `json:"iat,omitempty"`
	Issuer     string `json:"iss,omitempty"`
	Subject    string `json:"sub,omitempty"`
	Additional map[string]interface{}
	Scopes     []string `json:"scope,omitempty"`
}

func (claims *Claims) Serialize() jwt.MapClaims {
	mapClaims := jwt.MapClaims{
		"aud":   claims.Audience,
		"exp":   claims.ExpiresAt,
		"jti":   claims.Id,
		"iss":   claims.Issuer,
		"iat":   claims.IssuedAt,
		"sub":   claims.Subject,
		"scope": claims.Scopes,
	}

	for key, value := range claims.Additional {
		mapClaims[key] = value
	}
	return mapClaims
}

type ClaimParams struct {
	User          *users.UserDTO
	Application   *applications.Application
	ExpirationAdd int64
	Issuer        string
	Scopes        []string
	TokenType     string
}

const PasswordLogin = "pwd-login"

func makeClaims(ctx context.Context, params ClaimParams) Claims {
	iat := time.Now().Unix()
	id := params.TokenType + "." + uuid.NewV4().String()

	shared.GetLogger(ctx).WithFields(log.Fields{
		"id":  id,
		"iat": iat,
	}).Debug("Making token claims")

	audience := PasswordLogin
	if params.Application != nil {
		audience = params.Application.ClientId
	}

	return Claims{
		Audience:   audience,
		IssuedAt:   iat,
		ExpiresAt:  iat + params.ExpirationAdd,
		Id:         id,
		Issuer:     params.Issuer,
		Subject:    params.User.ID.String(),
		Scopes:     params.Scopes,
		Additional: make(map[string]interface{}),
	}
}
