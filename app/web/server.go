package web

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/infra"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

type WebServer struct {
	Router   *gin.Engine
	App      *infra.GouthyApp
	httpTool *web_utils.HTTPTools
}

func CreateWebServer(application *infra.GouthyApp) WebServer {
	router := gin.Default()

	return WebServer{
		Router:   router,
		App:      application,
		httpTool: web_utils.NewHTTPTools(application),
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

	newRestApi(s).Register(v1Route)
	registerWellKnownControllers(s.App, s.Router.Group("./.well-known"))
	newPages(s).Register(s.Router.Group("/pages"))

	return nil
}
