package apps

import "github.com/jinzhu/gorm"

type DiProvider struct {
	Repos    Repositories
	Services Services
	Facade	Facade
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

func newRepositories(db *gorm.DB) Repositories {
	return Repositories{
		Apps:   NweApplicationsRepositoryDB(db),
		Secrets: NewSecretsRepositoryDB(db),
	}
}

func NewDiProvider(db *gorm.DB) DiProvider {
	repos := newRepositories(db)
	services := newServices(repos.Apps)
	return DiProvider{
		Repos:    repos,
		Services: services,
		Facade: NewApplicationsFacade(repos.Apps, repos.Secrets, services.Find),
	}
}
