package core

import (
	"github.com/pestanko/gouthy/app/auth"
	"github.com/pestanko/gouthy/app/domain/apps"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/jwtlib"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/shared/store"
)

type GouthyApp struct {
	db      shared.DBConnection
	stores  store.Stores
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


// GetApplication - gets an application instance
func GetApplication(config *shared.AppConfig, connection shared.DBConnection) (*GouthyApp, error) {
	stores := store.NewRedisStoresFromConfig(&config.Redis)
	di := NewDI(connection, config, stores)
	return &GouthyApp{
		Config:  config,
		db:      connection,
		DI:      di,
		stores:  stores,
		Facades: newFacades(&di),
	}, nil
}

func NewDI(db shared.DBConnection, cfg *shared.AppConfig, stores store.Stores) DI {
	app := apps.NewDiProvider(db, &cfg.Features)
	user := users.NewDiProvider(db, &cfg.Features)
	jwtl := jwtlib.NewDiProvider(cfg.Jwk.Keys)
	authProvider := auth.NewDiProvider(
		app.Services.Find, user.Services.Find,
		jwtl.Services.Jwk, jwtl.Services.Jwt,
		user.Services.Password, stores)

	return DI{
		Auth:  authProvider,
		Users: user,
		Apps:  app,
	}
}

func newFacades(di *DI) Facades {
	return Facades{
		Auth:  di.Auth.Facade,
		Users: di.Users.Facade,
		Apps:  di.Apps.Facade,
		Keys:  di.Jwt.Facade,
	}
}
