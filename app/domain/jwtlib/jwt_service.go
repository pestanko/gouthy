package jwtlib

import (
	"context"
	"github.com/pestanko/gouthy/app/domain/apps"
	"github.com/pestanko/gouthy/app/domain/users"
)

type JwtService interface {
	CreateAccessToken(ctx context.Context, params TokenCreateParams) (Jwt, error)
	CreateRefreshToken(ctx context.Context, params TokenCreateParams) (Jwt, error)
	CreateIdToken(ctx context.Context, params TokenCreateParams) (Jwt, error)
	CreateSignedAccessToken(ctx context.Context, params TokenCreateParams) (*SignedJwt, error)
	CreateSignedRefreshToken(ctx context.Context, params TokenCreateParams) (*SignedJwt, error)
	CreateSignedIdToken(ctx context.Context, params TokenCreateParams) (*SignedJwt, error)
}



type jwtServiceImpl struct {
	keys       JwkRepository
	users      users.Repository
	apps       apps.Repository
	jwtSigning JwtSigningService
}

func (j *jwtServiceImpl) CreateSignedAccessToken(ctx context.Context, params TokenCreateParams) (*SignedJwt, error) {
	token, err := j.CreateAccessToken(ctx, params)
	if err != nil {
		return nil, err
	}

	return j.signToken(ctx, token)
}

func (j *jwtServiceImpl) CreateSignedRefreshToken(ctx context.Context, params TokenCreateParams) (*SignedJwt, error) {
	token, err := j.CreateRefreshToken(ctx, params)
	if err != nil {
		return nil, err
	}

	return j.signToken(ctx, token)
}

func (j *jwtServiceImpl) CreateSignedIdToken(ctx context.Context, params TokenCreateParams) (*SignedJwt, error) {
	token, err := j.CreateIdToken(ctx, params)
	if err != nil {
		return nil, err
	}

	return j.signToken(ctx, token)
}

func (j *jwtServiceImpl) CreateAccessToken(ctx context.Context, params TokenCreateParams) (Jwt, error) {
	const ExpTime = HOUR
	claims, err := j.makeClaims(ctx, params, ExpTime, "access_token")
	if err != nil {
		return nil, err
	}

	return j.createToken(ctx, claims)
}

func (j *jwtServiceImpl) CreateRefreshToken(ctx context.Context, params TokenCreateParams) (Jwt, error) {
	const ExpTime = HOUR * 24 * 7
	claims, err := j.makeClaims(ctx, params, ExpTime, "refresh_token")
	if err != nil {
		return nil, err
	}

	return j.createToken(ctx, claims)
}

func (j *jwtServiceImpl) CreateIdToken(ctx context.Context, params TokenCreateParams) (Jwt, error) {
	const ExpTime = HOUR * 8
	claims, err := j.makeClaims(ctx, params, ExpTime, "id_token")
	if err != nil {
		return nil, err
	}

	return j.createToken(ctx, claims)
}

func (j *jwtServiceImpl) createToken(ctx context.Context, claims Claims) (Jwt, error) {
	jwt, err := j.jwtSigning.Create(ctx, claims)
	if err != nil {
		return nil, err
	}
	return jwt, nil
}

func (j *jwtServiceImpl) signToken(ctx context.Context, jwt Jwt) (*SignedJwt, error) {
	signed, err := j.jwtSigning.Sign(ctx, jwt)

	if err != nil {
		return nil, err
	}
	return signed, nil
}

func (j *jwtServiceImpl) makeClaims(ctx context.Context, params TokenCreateParams, expTime int64, tokenType string) (Claims, error) {
	claims := makeClaims(ctx, ClaimParams{
		User:          params.User,
		Application:   params.App,
		Scopes:        params.Scopes,
		ExpirationAdd: expTime,
		TokenType:     tokenType,
		Issuer:        "Gouthy", // TODO
	})
	return claims, nil
}

func NewJwtService(keys JwkRepository, users users.Repository, apps apps.Repository) JwtService {
	jwtSigning := NewJwtSigningService(keys)
	return &jwtServiceImpl{
		keys:       keys,
		users:      users,
		apps:       apps,
		jwtSigning: jwtSigning,
	}
}
