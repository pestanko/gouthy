package web

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/app/web/api"
)

type WebServer struct {
	Router *gin.Engine
	App    *core.GouthyApp
}

func CreateWebServer(application *core.GouthyApp) WebServer {
	router := gin.Default()

	return WebServer{
		Router: router,
		App:    application,
	}
}

func (s *WebServer) Serve() error {
	return s.Router.Run(s.App.Config.Server.Port)
}

func (s *WebServer) Run() error {

	if err := s.RegisterRoutes(); err != nil {
		return err
	}

	return s.Serve()
}

func (s *WebServer) RegisterRoutes() error {
	apiRoute := s.Router.Group("/api")
	v1Route := apiRoute.Group("/v1")
	api.RegisterApiControllers(s.App, v1Route)

	wellKnownRoute := s.Router.Group("./.well-known")
	RegisterWellKnownControllers(s.App, wellKnownRoute)

	return nil
}

