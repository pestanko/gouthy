package infra

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pestanko/gouthy/app/domain/applications"
	"github.com/pestanko/gouthy/app/domain/auth"
	"github.com/pestanko/gouthy/app/domain/jwtlib"
	"github.com/pestanko/gouthy/app/domain/users"
)

type GouthyApp struct {
	db      *gorm.DB
	Config  *AppConfig
	Facades Facades
}

type Repositories struct {
	Users       users.Repository
	UserSecrets users.SecretsRepository
	Jwk         jwtlib.JwkRepository
	Apps        applications.Repository
	AppsSecrets applications.SecretsRepository
}

type Facades struct {
	Auth  auth.Facade
	Users users.Facade
	Apps  applications.Facade
	Keys  jwtlib.KeysFacade
}

func GetDBConnection(config *AppConfig) (*gorm.DB, error) {
	return gorm.Open(
		"postgres",
		fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
			config.DB.Host, config.DB.Port, config.DB.User, config.DB.DBName, config.DB.Password, config.DB.SSLMode),
	)
}

// GetApplication - gets an application instance
func GetApplication(config *AppConfig, db *gorm.DB) (GouthyApp, error) {
	app := GouthyApp{Config: config, db: db}
	registerFacades := NewFacades(&app)
	app.Facades = registerFacades
	return app, nil
}

func NewRepositories(app *GouthyApp) Repositories {
	return Repositories{
		Users:       users.NewUsersRepositoryDB(app.db),
		UserSecrets: users.NewSecretsRepositoryDB(app.db),
		Apps:        applications.NweApplicationsRepositoryDB(app.db),
		AppsSecrets: applications.NewSecretsRepositoryDB(app.db),
		Jwk:         jwtlib.NewJwkRepository(app.Config.Jwk.Keys),
	}
}

func NewFacades(app *GouthyApp) Facades {
	repos := NewRepositories(app)

	return Facades{
		Auth:  auth.NewAuthFacade(repos.Users, repos.Apps, repos.Jwk),
		Users: users.NewUsersFacade(repos.Users, repos.UserSecrets),
		Apps:  applications.NewApplicationsFacade(repos.Apps, repos.AppsSecrets),
		Keys:  jwtlib.NewKeysFacade(repos.Users, repos.Jwk),
	}
}
