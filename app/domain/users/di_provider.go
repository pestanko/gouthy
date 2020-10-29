package users

import (
	"github.com/pestanko/gouthy/app/shared"
)

type DiProvider struct {
	Repos    Repositories
	Services Services
	Facade   Facade
}

type Repositories struct {
	Users   Repository
	Secrets SecretsRepository
}

type Services struct {
	Find     FindService
	Password PasswordService
}

func newServices(usersRepo Repository) Services {
	return Services{
		Find:     NewFindUsersService(usersRepo),
		Password: NewPasswordService(usersRepo),
	}
}

func newRepositories(connection shared.DBConnection) Repositories {
	gorm := shared.DBConnectionIntoGorm(connection)
	return Repositories{
		Users:   NewUsersRepositoryDB(gorm),
		Secrets: NewSecretsRepositoryDB(gorm),
	}
}

func NewDiProvider(db shared.DBConnection, features *shared.FeaturesConfig) DiProvider {
	repos := newRepositories(db)
	services := newServices(repos.Users)
	return DiProvider{
		Repos:    repos,
		Services: services,
		Facade:   NewUsersFacade(repos.Users, repos.Secrets, services, features),
	}
}
