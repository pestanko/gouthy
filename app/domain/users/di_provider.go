package users

import "github.com/jinzhu/gorm"

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
	Find FindService
	Password PasswordService
}

func newServices(usersRepo Repository) Services {
	return Services{
		Find:     NewFindUsersService(usersRepo),
		Password: NewPasswordService(usersRepo),
	}
}

func newRepositories(db *gorm.DB) Repositories {
	return Repositories{
		Users:   NewUsersRepositoryDB(db),
		Secrets: NewSecretsRepositoryDB(db),
	}
}

func NewDiProvider(db *gorm.DB) DiProvider {
	repos := newRepositories(db)
	services := newServices(repos.Users)
	return DiProvider{
		Repos:    repos,
		Services: services,
		Facade:   NewUsersFacade(repos.Users, repos.Secrets, services),
	}
}