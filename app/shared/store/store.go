package store

import (
	"context"
	"time"
)

const (
	StoreUnspecifiedDB		   = 0
	StoreActivationCodeDB      = 1
	StoreOAuth2AuthorizationDB = 2
)

type Stores interface {
	GetStore(id int) Store
}

type Store interface {
	Get(ctx context.Context, key string) (string, error)
	GetJson(ctx context.Context, key string, value interface{}) error
	Set(ctx context.Context, key string, value interface{}, exp time.Duration) error
	SetJson(ctx context.Context, key string, value interface{}, exp time.Duration) error
	Delete(ctx context.Context, key string) error
	ListKeys(ctx context.Context, pattern string) ([]string, error)
}
