package jwtlib

import (
	"context"
)

type KeysFacade interface {
	ListJwks(ctx context.Context) ([]Jwk, error)
	GenerateNewJwk(ctx context.Context) error
	GetLatest(ctx context.Context) (Jwk, error)
}

type keysFacadeImpl struct {
	jwk JwkService
}

func (keys *keysFacadeImpl) GetLatest(ctx context.Context) (Jwk, error) {
	return keys.GetLatest(ctx)
}

func (keys *keysFacadeImpl) ListJwks(ctx context.Context) ([]Jwk, error) {
	return keys.jwk.List(ctx)
}

func (keys *keysFacadeImpl) GenerateNewJwk(ctx context.Context) error {
	return keys.jwk.GenerateNew(ctx)
}

func NewKeysFacade(jwkService JwkService) KeysFacade {
	return &keysFacadeImpl{jwk: jwkService}
}
