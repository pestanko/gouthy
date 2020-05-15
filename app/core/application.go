package core

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pestanko/gouthy/app/domain/auth"
	"github.com/pestanko/gouthy/app/domain/auth/jwtlib"
	"github.com/pestanko/gouthy/app/domain/entities"
	"github.com/pestanko/gouthy/app/domain/users"
)

type GouthyApp struct {
	Config  *AppConfig
	DB      *gorm.DB
	Facades Facades
}

type Facades struct {
	Auth     auth.Facade
	Users    users.Facade
	Entities entities.Facade
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
	app := GouthyApp{Config: config, DB: db}
	registerFacades := RegisterFacades(&app)
	app.Facades = registerFacades
	return app, nil
}

func RegisterFacades(app *GouthyApp) Facades {
	jwkInventory := jwtlib.NewJwkInventory(app.Config.Jwk.Keys)
	usersFacade := users.NewUsersFacade(app.DB)
	entitiesFacade := entities.NewEntitiesFacade(app.DB)
	authFacade := auth.NewAuthFacade(app.DB, usersFacade, entitiesFacade, jwkInventory)

	return Facades{
		Auth:     authFacade,
		Users:    usersFacade,
		Entities: entitiesFacade,
	}
}
