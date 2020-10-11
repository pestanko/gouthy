package jwtlib

import (
	"context"
)

const (
	TokenTypeAccess  = "a"
	TokenTypeRefresh = "r"
	TokenTypeId      = "i"
	TokenTypeSession = "s"
	TokenTypeUndefined = "U"
)

type JwtService interface {
	CreateAccessToken(ctx context.Context, params TokenCreateParams) (Jwt, error)
	CreateRefreshToken(ctx context.Context, params TokenCreateParams) (Jwt, error)
	CreateIdToken(ctx context.Context, params TokenCreateParams) (Jwt, error)
	CreateSessionToken(ctx context.Context, params TokenCreateParams) (Jwt, error)
	CreateSignedAccessToken(ctx context.Context, params TokenCreateParams) (*SignedJwt, error)
	CreateSignedRefreshToken(ctx context.Context, params TokenCreateParams) (*SignedJwt, error)
	CreateSignedIdToken(ctx context.Context, params TokenCreateParams) (*SignedJwt, error)
	CreateSignedSessionToken(ctx context.Context, params TokenCreateParams) (*SignedJwt, error)
}

func NewJwtService(keys JwkRepository) JwtService {
	jwtSigning := NewJwtSigningService(keys)
	return &jwtServiceImpl{
		keys:       keys,
		jwtSigning: jwtSigning,
	}
}

type jwtServiceImpl struct {
	keys       JwkRepository
	jwtSigning JwtSigningService
}

func (j *jwtServiceImpl) CreateSessionToken(ctx context.Context, params TokenCreateParams) (Jwt, error) {
	const ExpTime = SessionTokenExpiration
	return j.createInternal(ctx, params, ExpTime, TokenTypeRefresh)
}

func (j *jwtServiceImpl) CreateSignedSessionToken(ctx context.Context, params TokenCreateParams) (*SignedJwt, error) {
	token, err := j.CreateSessionToken(ctx, params)
	if err != nil {
		return nil, err
	}
	return j.signToken(ctx, token)
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
	const ExpTime = AccessTokenExpiration
	return j.createInternal(ctx, params, ExpTime, TokenTypeAccess)
}

func (j *jwtServiceImpl) CreateRefreshToken(ctx context.Context, params TokenCreateParams) (Jwt, error) {
	const ExpTime = RefreshTokenExpiration
	return j.createInternal(ctx, params, ExpTime, TokenTypeRefresh)
}

func (j *jwtServiceImpl) CreateIdToken(ctx context.Context, params TokenCreateParams) (Jwt, error) {
	const ExpTime = IdTokenExpiration
	return j.createInternal(ctx, params, ExpTime, TokenTypeId)
}

func (j *jwtServiceImpl) createInternal(ctx context.Context, params TokenCreateParams, exp int64, tokenType string) (Jwt, error) {
	claims, err := j.makeClaims(ctx, params, exp, tokenType)
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
	claims := makeClaims(ctx, claimsParams{
		User:          params.User,
		Application:   params.App,
		Scopes:        params.Scopes,
		ExpirationAdd: expTime,
		TokenType:     tokenType,
		Issuer:        "Gouthy", // TODO
	})
	return claims, nil
}
