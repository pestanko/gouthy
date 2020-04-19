package core

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pestanko/gouthy/app/jwtl"
	"github.com/pestanko/gouthy/app/services"
)

type GouthyApp struct {
	Config   *AppConfig
	DB       *gorm.DB
	Services Services
}

type Services struct {
	Auth     *services.AuthService
	Users    services.UsersService
	Entities services.EntitiesService
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
	registerServices := RegisterServices(&app)
	app.Services = registerServices
	return app, nil
}

func RegisterServices(app *GouthyApp) Services {
	jwkInventory := jwtl.NewJwkInventory(app.Config.Jwk.Keys)
	users := services.NewUsersService(app.DB)
	entities := services.NewEntitiesService(app.DB)
	auth := services.NewAuthService(app.DB, users, entities, jwkInventory)

	return Services{
		Auth:     auth,
		Users:    users,
		Entities: entities,
	}
}
