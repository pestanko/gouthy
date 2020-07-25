package jwtlib

type DiProvider struct {
	Repos    Repositories
	Services Services
	Facade   KeysFacade
}

type Repositories struct {
	Jwk JwkRepository
}

type Services struct {
	Jwk JwkService
	Jwt JwtService
}

func newServices(jwkRepo JwkRepository) Services {
	return Services{
		Jwk: NewJwkService(jwkRepo),
		Jwt: NewJwtService(jwkRepo),
	}
}

func newRepositories(keyPath string) Repositories {
	return Repositories{
		Jwk: NewJwkRepository(keyPath),
	}
}

func NewDiProvider(keyPath string) DiProvider {
	repos := newRepositories(keyPath)
	services := newServices(repos.Jwk)
	return DiProvider{
		Repos:    repos,
		Services: services,
		Facade: NewKeysFacade(services.Jwk),
	}
}
