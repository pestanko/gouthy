package auth

import (
	"github.com/pestanko/gouthy/app/domain/apps"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/jwtlib"
	"github.com/pestanko/gouthy/app/shared/store"
)

type DiProvider struct {
	Repos    Repos
	Services Services
	Facade   Facade
}

type Services struct {
	OAuth2Service OAuth2AuthorizationService
}

type Repos struct {
	OAuth2AuthReq OAuth2AuthRequestRepo
}

func newServices(appFind apps.FindService) Services {
	return Services{
		OAuth2Service: NewOAuth2AuthorizationService(appFind),
	}
}

func newRepos(stores store.Stores) Repos {
	return Repos{
		OAuth2AuthReq: NewOAuth2AuthRequestRepo(stores.GetStore(store.StoreOAuth2AuthorizationDB)),
	}
}

func NewDiProvider(
	appFind apps.FindService,
	userFind users.FindService,
	jwk jwtlib.JwkService,
	jwt jwtlib.JwtService,
	pwdService users.PasswordService,
	stores store.Stores) DiProvider {

	services := newServices(appFind)
	return DiProvider{
		Repos:    newRepos(stores),
		Services: services,
		Facade:   NewAuthFacade(userFind, appFind, jwt, jwk, pwdService),
	}
}
