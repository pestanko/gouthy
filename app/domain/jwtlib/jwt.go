package jwtlib

import (
	"context"
	"github.com/dgrijalva/jwt-go"
)

type Jwt interface {
}

type JWTImpl struct {
	token *jwt.Token
}

func NewJwt(token *jwt.Token) Jwt {
	return &JWTImpl{
		token: token,
	}
}

type JwtSigningService interface {
	Sign(ctx context.Context, token Jwt) (Jwt, error)
	Create(ctx context.Context, claims jwt.Claims) (Jwt, error)
}

type jwtSigningServiceImpl struct {
	repo JwkRepository
}

func (service *jwtSigningServiceImpl) Sign(ctx context.Context, token Jwt) (Jwt, error) {
	return token, nil
}

func (service *jwtSigningServiceImpl) Create(ctx context.Context, claims jwt.Claims) (Jwt, error) {
	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
	return NewJwt(token), nil
}

func NewJwtSigningService(repo JwkRepository) JwtSigningService {
	return &jwtSigningServiceImpl{repo: repo}
}
