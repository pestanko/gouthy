package jwtlib

import (
	"context"
	"github.com/pestanko/gouthy/app/domain/applications"
	"github.com/pestanko/gouthy/app/domain/users"
)

const HOUR = 3600

type JwkService interface {
	GenerateNew(ctx context.Context) error
	List(ctx context.Context) ([]Jwk, error)
	Get(ctx context.Context, id string) (Jwk, error)
	GetLatest(ctx context.Context) (Jwk, error)
}

type jwkServiceImpl struct {
	repo      JwkRepository
	usersRepo users.Repository
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

func NewJwkService(jwkRepo JwkRepository, usersRepo users.Repository) JwkService {
	return &jwkServiceImpl{repo: jwkRepo, usersRepo: usersRepo}
}

type TokenCreateParams struct {
	User   *users.UserDTO
	App    *applications.Application
	Scopes []string
}

