package infra

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pestanko/gouthy/app/domain/apps"
	"github.com/pestanko/gouthy/app/domain/auth"
	"github.com/pestanko/gouthy/app/domain/jwtlib"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/shared"
)

type GouthyApp struct {
	db      *gorm.DB
	Config  *shared.AppConfig
	Facades Facades
	DI      DI
}

type Facades struct {
	Auth  auth.Facade
	Users users.Facade
	Apps  apps.Facade
	Keys  jwtlib.KeysFacade
}

type DI struct {
	Auth  auth.DiProvider
	Users users.DiProvider
	Apps  apps.DiProvider
	Jwt   jwtlib.DiProvider
}

func GetDBConnection(config *shared.AppConfig) (*gorm.DB, error) {
	return gorm.Open(
		"postgres",
		fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
			config.DB.Host, config.DB.Port, config.DB.User, config.DB.DBName, config.DB.Password, config.DB.SSLMode),
	)
}

// GetApplication - gets an application instance
func GetApplication(config *shared.AppConfig, db *gorm.DB) (GouthyApp, error) {
	app := GouthyApp{Config: config, db: db, DI: NewDI(db, config)}
	app.Facades = newFacades(&app)
	return app, nil
}

func NewDI(db *gorm.DB, cfg *shared.AppConfig) DI {
	app := apps.NewDiProvider(db)
	user := users.NewDiProvider(db)
	jwtl := jwtlib.NewDiProvider(cfg.Jwk.Keys)
	authProvider := auth.NewDiProvider(app.Services.Find, user.Services.Find, jwtl.Services.Jwk, jwtl.Services.Jwt, user.Services.Password)
	return DI{
		Auth:  authProvider,
		Users: user,
		Apps:  app,
	}
}

func newFacades(app *GouthyApp) Facades {
	return Facades{
		Auth:  app.DI.Auth.Facade,
		Users: app.DI.Users.Facade,
		Apps:  app.DI.Apps.Facade,
		Keys:  app.DI.Jwt.Facade,
	}
}
