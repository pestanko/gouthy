package web

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/infra"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

type GouthyWebServer struct {
	Router   *gin.Engine
	App      *infra.GouthyApp
	httpTool *web_utils.HTTPTools
	controllers Controllers
}

func CreateWebServer(application *infra.GouthyApp) *GouthyWebServer {
	tools := web_utils.NewHTTPTools(application)
	return &GouthyWebServer{
		Router:      gin.Default(),
		App:         application,
		httpTool:    tools,
		controllers: CreateControllers(tools),
	}
}

func (s *GouthyWebServer) Serve() error {
	return s.Router.Run(s.App.Config.Server.Port)
}

func (s *GouthyWebServer) Run() error {
	RegisterRoutes(s)
	return s.Serve()
}


type Controllers struct {
	pages *pagesControllers
	api   *apisControllers
	WellKnown *WellKnownController
}


func CreateControllers(tools *web_utils.HTTPTools) Controllers {
	return Controllers{
		pages:     NewPagesControllers(tools),
		api:       NewApiControllers(tools),
		WellKnown: NewWellKnownController(tools),
	}
}


