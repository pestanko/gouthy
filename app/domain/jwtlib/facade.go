package jwtlib

import (
	"context"
	"github.com/pestanko/gouthy/app/domain/applications"
	"github.com/pestanko/gouthy/app/domain/users"
)

const HOUR = 3600

type JwkFacade interface {
	GenerateNew(ctx context.Context) error
	List(ctx context.Context) ([]Jwk, error)
	Get(ctx context.Context, id string) (Jwk, error)
	GetLatest(ctx context.Context) (Jwk, error)
}

type jwkFacadeImpl struct {
	repo      JwkRepository
	usersRepo users.Repository
}

func (facade *jwkFacadeImpl) GetLatest(ctx context.Context) (Jwk, error) {
	return facade.repo.GetLatest(ctx)
}

func (facade *jwkFacadeImpl) Get(ctx context.Context, id string) (Jwk, error) {
	return facade.repo.Get(ctx, id)
}

func (facade *jwkFacadeImpl) List(ctx context.Context) ([]Jwk, error) {
	return facade.repo.List(ctx)
}

func (facade *jwkFacadeImpl) GenerateNew(ctx context.Context) error {
	return facade.repo.Generate(ctx, JwkGenerateParams{})
}

func NewJwkFacade(jwkRepo JwkRepository, usersRepo users.Repository) JwkFacade {
	return &jwkFacadeImpl{repo: jwkRepo, usersRepo: usersRepo}
}

type TokenCreateParams struct {
	User   *users.UserDTO
	App    *applications.ApplicationDTO
	Scopes []string
}

type JwtFacade interface {
	CreateAccessToken(ctx context.Context, params TokenCreateParams) (Jwt, error)
	CreateRefreshToken(ctx context.Context, params TokenCreateParams) (Jwt, error)
	CreateIdToken(ctx context.Context, params TokenCreateParams) (Jwt, error)
	CreateSignedAccessToken(ctx context.Context, params TokenCreateParams) (*SignedJwt, error)
	CreateSignedRefreshToken(ctx context.Context, params TokenCreateParams) (*SignedJwt, error)
	CreateSignedIdToken(ctx context.Context, params TokenCreateParams) (*SignedJwt, error)
}

type JwtFacadeImpl struct {
	keys       JwkRepository
	users      users.Repository
	apps       applications.Repository
	jwtSigning JwtSigningService
}

func (j *JwtFacadeImpl) CreateSignedAccessToken(ctx context.Context, params TokenCreateParams) (*SignedJwt, error) {
	token, err := j.CreateAccessToken(ctx, params)
	if err != nil {
		return nil, err
	}

	return j.signToken(ctx, token)
}

func (j *JwtFacadeImpl) CreateSignedRefreshToken(ctx context.Context, params TokenCreateParams) (*SignedJwt, error) {
	token, err := j.CreateRefreshToken(ctx, params)
	if err != nil {
		return nil, err
	}

	return j.signToken(ctx, token)
}

func (j *JwtFacadeImpl) CreateSignedIdToken(ctx context.Context, params TokenCreateParams) (*SignedJwt, error) {
	token, err := j.CreateIdToken(ctx, params)
	if err != nil {
		return nil, err
	}

	return j.signToken(ctx, token)
}

func (j *JwtFacadeImpl) CreateAccessToken(ctx context.Context, params TokenCreateParams) (Jwt, error) {
	const ExpTime = HOUR
	claims, err := j.makeClaims(ctx, params, ExpTime, "access_token")
	if err != nil {
		return nil, err
	}

	return j.createToken(ctx, claims)
}

func (j *JwtFacadeImpl) CreateRefreshToken(ctx context.Context, params TokenCreateParams) (Jwt, error) {
	const ExpTime = HOUR * 24 * 7
	claims, err := j.makeClaims(ctx, params, ExpTime, "refresh_token")
	if err != nil {
		return nil, err
	}

	return j.createToken(ctx, claims)
}

func (j *JwtFacadeImpl) CreateIdToken(ctx context.Context, params TokenCreateParams) (Jwt, error) {
	const ExpTime = HOUR * 8
	claims, err := j.makeClaims(ctx, params, ExpTime, "id_token")
	if err != nil {
		return nil, err
	}

	return j.createToken(ctx, claims)
}

func (j *JwtFacadeImpl) createToken(ctx context.Context, claims Claims) (Jwt, error) {
	jwt, err := j.jwtSigning.Create(ctx, claims)
	if err != nil {
		return nil, err
	}
	return jwt, nil
}

func (j *JwtFacadeImpl) signToken(ctx context.Context, jwt Jwt) (*SignedJwt, error) {
	signed, err := j.jwtSigning.Sign(ctx, jwt)

	if err != nil {
		return nil, err
	}
	return signed, nil
}

func (j *JwtFacadeImpl) makeClaims(ctx context.Context, params TokenCreateParams, expTime int64, tokenType string) (Claims, error) {
	claims := MakeClaims(ctx, ClaimParams{
		User:          params.User,
		Application:   params.App,
		Scopes:        params.Scopes,
		ExpirationAdd: expTime,
		TokenType:     tokenType,
		Issuer:        "Gouthy", // TODO
	})
	return claims, nil
}

func NewJwtFacade(keys JwkRepository, users users.Repository, apps applications.Repository) JwtFacade {
	jwtSigning := NewJwtSigningService(keys)
	return &JwtFacadeImpl{
		keys:       keys,
		users:      users,
		apps:       apps,
		jwtSigning: jwtSigning,
	}
}
