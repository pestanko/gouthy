package jwtlib

import (
	"context"
	"github.com/pestanko/gouthy/app/domain/apps"
	"github.com/pestanko/gouthy/app/domain/users"
)

const HOUR int64 = 3600
const DAY = HOUR * 24
const AccessTokenExpiration = HOUR
const RefreshTokenExpiration = 7 * DAY // WEEK
const IdTokenExpiration = 8 * HOUR
const SessionTokenExpiration = 12 * HOUR

type JwkService interface {
	GenerateNew(ctx context.Context) error
	List(ctx context.Context) ([]Jwk, error)
	Get(ctx context.Context, id string) (Jwk, error)
	GetLatest(ctx context.Context) (Jwk, error)
}

type jwkServiceImpl struct {
	repo JwkRepository
}

func (facade *jwkServiceImpl) GetLatest(ctx context.Context) (Jwk, error) {
	return facade.repo.GetLatest(ctx)
}

func (facade *jwkServiceImpl) Get(ctx context.Context, id string) (Jwk, error) {
	return facade.repo.Get(ctx, id)
}

func (facade *jwkServiceImpl) List(ctx context.Context) ([]Jwk, error) {
	return facade.repo.List(ctx)
}

func (facade *jwkServiceImpl) GenerateNew(ctx context.Context) error {
	return facade.repo.Generate(ctx, JwkGenerateParams{})
}

func NewJwkService(jwkRepo JwkRepository) JwkService {
	return &jwkServiceImpl{repo: jwkRepo}
}

type TokenCreateParams struct {
	User          *users.User
	App           *apps.Application
	Scopes        []string
	CorrelationId string
}
