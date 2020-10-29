package apps

import (
	"github.com/pestanko/gouthy/app/shared"
)

type DiProvider struct {
	Repos    Repositories
	Services Services
	Facade   Facade
}

type Repositories struct {
	Apps    Repository
	Secrets SecretsRepository
}

type Services struct {
	Find FindService
}

func newServices(appsRepo Repository) Services {
	return Services{
		Find: NewFindAppsService(appsRepo),
	}
}

func newRepositories(connection shared.DBConnection) Repositories {
	gorm := shared.DBConnectionIntoGorm(connection)
	return Repositories{
		Apps:    NweApplicationsRepositoryDB(gorm),
		Secrets: NewSecretsRepositoryDB(gorm),
	}
}

func NewDiProvider(db shared.DBConnection, features *shared.FeaturesConfig) DiProvider {
	repos := newRepositories(db)
	services := newServices(repos.Apps)
	return DiProvider{
		Repos:    repos,
		Services: services,
		Facade:   NewApplicationsFacade(repos.Apps, repos.Secrets, services.Find),
	}
}
