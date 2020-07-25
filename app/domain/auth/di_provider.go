package auth

import (
	"github.com/pestanko/gouthy/app/domain/apps"
	"github.com/pestanko/gouthy/app/domain/jwtlib"
	"github.com/pestanko/gouthy/app/domain/users"
)

type DiProvider struct {
	Services Services
	Facade Facade
}

type Services struct {
	OAuth2Service OAuth2AuthorizationService
}

func newServices(appFind apps.FindService) Services {
	return Services{
		OAuth2Service: NewOAuth2AuthorizationService(appFind),
	}
}


func NewDiProvider(appFind apps.FindService, userFind users.FindService, jwk jwtlib.JwkService, jwt jwtlib.JwtService, passwdService users.PasswordService) DiProvider {
	services := newServices(appFind)
	return DiProvider{
		Services: services,
		Facade:  NewAuthFacade(userFind, appFind, jwt, jwk, passwdService),
	}
}
