package auth

import (
	"context"
	"github.com/pestanko/gouthy/app/shared/store"
	"time"
)

const OAuth2AuthRequestDuration = 60 * time.Second

// https://tools.ietf.org/html/rfc6749#section-4.1.1
type OAuth2AuthRequest struct {
	ClientId      string   `json:"client_id"`
	ResponseType  string   `json:"response_type"`
	RedirectUri   string   `json:"redirect_uri"`
	Scopes        []string `json:"scopes"`
	State         string   `json:"state"`
	Nonce         string   `json:"nonce"`               // OPENID
	PKCEChallenge string   `json:"pkce_code_challenge"` // PKCE
	PKCEMethod    string   `json:"pkce_code_method"`    // PKCE
}

type OAuth2AuthRequestRepo interface {
	Get(ctx context.Context, key string) (OAuth2AuthRequest, error)
	Set(ctx context.Context, key string, request OAuth2AuthRequest) error
	Delete(ctx context.Context, key string) error
}

func NewOAuth2AuthRequestRepo(s store.Store) OAuth2AuthRequestRepo {
	return &oAuth2AuthRequestRepo{
		s: s,
	}
}

type oAuth2AuthRequestRepo struct {
	s store.Store
}

func (repo *oAuth2AuthRequestRepo) Get(ctx context.Context, key string) (OAuth2AuthRequest, error) {
	var request OAuth2AuthRequest
	if err := repo.s.GetJson(ctx, key, &request); err != nil {
		return OAuth2AuthRequest{}, err
	}
	return request, nil
}

func (repo *oAuth2AuthRequestRepo) Set(ctx context.Context, key string, request OAuth2AuthRequest) error {
	return repo.s.SetJson(ctx, key, &request, OAuth2AuthRequestDuration)
}

func (repo *oAuth2AuthRequestRepo) Delete(ctx context.Context, key string) error {
	return repo.s.Delete(ctx, key)
}
