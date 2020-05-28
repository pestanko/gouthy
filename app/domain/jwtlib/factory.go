package jwtlib

import "context"

type JwkFacade interface {
	GenerateNew(ctx context.Context) error
	List(ctx context.Context) ([]Jwk, error)
	Get(ctx context.Context, id string) (Jwk, error)
}

type jwkFacadeImpl struct {
	repo JwkRepository
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

func NewJwkFacade(jwkRepo JwkRepository) JwkFacade {
	return &jwkFacadeImpl{repo: jwkRepo}
}
