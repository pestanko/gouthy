package web

import (
	"github.com/gin-gonic/gin"
	ctrl "github.com/pestanko/gouthy/app/web/controllers"
	"github.com/pestanko/gouthy/app/web/shared"
)

type RestApi struct {
	server *WebServer
}

func NewRestApi(server *WebServer) *RestApi {
	return &RestApi{server: server}
}

func (api *RestApi) Register(r *gin.RouterGroup) []shared.Controller {
	authController := api.NewAuthController()
	usersController := api.NewUsersController()
	machinesController := api.NewMachinesController()
	entityController := api.NewEntityController()
	var ctrls = []shared.Controller{
		authController.RegisterRoutes(r.Group("/auth")),
		usersController.RegisterRoutes(r.Group("/users")),
		machinesController.RegisterRoutes(r.Group("/machines")),
		entityController.RegisterRoutes(r.Group("/entities")),
	}

	return ctrls
}

func (api *RestApi) NewUsersController() *ctrl.UsersController {
	return &ctrl.UsersController{
		Users: api.server.App.Facades.Users,
		Http:  api.server.httpTool,
	}
}

func (api *RestApi) NewEntityController() *ctrl.EntitiesController {
	return &ctrl.EntitiesController{
		Entities: api.server.App.Facades.Entities,
		Http:     api.server.httpTool,
	}
}

func (api *RestApi) NewMachinesController() *ctrl.MachinesController {
	return &ctrl.MachinesController{Machines: nil, Http: api.server.httpTool}
}

func (api *RestApi) NewAuthController() *ctrl.AuthController {
	entitiesFacade := api.server.App.Facades.Entities
	usersFacade := api.server.App.Facades.Users
	return &ctrl.AuthController{
		Entities: entitiesFacade,
		Users:    usersFacade,
		Auth:     api.server.App.Facades.Auth,
		Http:     api.server.httpTool,
	}
}
