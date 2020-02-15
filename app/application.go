package app

import (
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/services"
)

type Application struct {
	DB       gorm.DB
	Accounts services.AccountsService
	Users    services.UsersService
	Secrets  services.SecretsService
}

func NewApplication(db gorm.DB) Application {
	return Application{
		DB:       db,
		Accounts: services.NewAccountsService(db),
		Users:    services.NewUsersService(db),
		Secrets:  services.NewSecretsService(db),
	}
}
