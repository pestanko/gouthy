package web

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

type AppWebServer struct {
	Router      *gin.Engine
	App         *core.GouthyApp
	httpTool    *web_utils.Tools
	controllers Controllers
}

func CreateWebServer(application *core.GouthyApp) *AppWebServer {
	tools := web_utils.NewHTTPTools(application)
	return &AppWebServer{
		Router:      gin.Default(),
		App:         application,
		httpTool:    tools,
		controllers: CreateControllers(tools),
	}
}

func (s *AppWebServer) Serve() error {
	return s.Router.Run(s.App.Config.Server.Port)
}

func (s *AppWebServer) Run() error {
	RegisterRoutes(s)
	return s.Serve()
}

type Controllers struct {
	pages     *pagesControllers
	api       *apisControllers
	WellKnown *WellKnownController
}

func CreateControllers(tools *web_utils.Tools) Controllers {
	return Controllers{
		pages:     NewPagesControllers(tools),
		api:       NewApiControllers(tools),
		WellKnown: NewWellKnownController(tools),
	}
}
