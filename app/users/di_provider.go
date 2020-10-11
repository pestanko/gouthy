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

func newRepositories(db shared.DBConnection) Repositories {
	gorm := shared.DBConnectionIntoGorm(db)
	return Repositories{
		Users:   NewUsersRepositoryDB(gorm),
		Secrets: NewSecretsRepositoryDB(gorm),
	}
}

func NewDiProvider(db shared.DBConnection) DiProvider {
	repos := newRepositories(db)
	services := newServices(repos.Users)
	return DiProvider{
		Repos:    repos,
		Services: services,
		Facade:   NewUsersFacade(repos.Users, repos.Secrets, services),
	}
}