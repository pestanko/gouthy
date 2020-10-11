package jwtlib

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/pestanko/gouthy/app/apps"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/users"
	log "github.com/sirupsen/logrus"
	"strconv"
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
	Scopes     shared.Scopes `json:"scope,omitempty"`
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

type claimsParams struct {
	User          *users.User
	Application   *apps.Application
	ExpirationAdd int64
	Issuer        string
	Scopes        shared.Scopes
	TokenType     string
	CorrelationId string
	Iat           int64
}

const PasswordLogin = "pwd-login"

func makeClaims(ctx context.Context, params claimsParams) Claims {
	iat := time.Now().Unix()
	id := JtiPartsToString(JtiParts{
		CorrelationId: params.CorrelationId,
		Type:          params.TokenType,
		Offset:        makeOffset(params),
	})

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

func makeOffset(params claimsParams) string {
	switch params.TokenType {
	case TokenTypeAccess:
		return strconv.FormatInt(params.Iat, 10)
	case TokenTypeRefresh:
		return "0"
	case TokenTypeId:
		return "0"
	case TokenTypeSession:
		return "0"
	}
	return TokenTypeUndefined
}
