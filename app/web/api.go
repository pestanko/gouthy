package web

import (
	"github.com/gin-gonic/gin"
	ctrl "github.com/pestanko/gouthy/app/web/api"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

type restApi struct {
	server *WebServer
}

func newRestApi(server *WebServer) *restApi {
	return &restApi{server: server}
}

func (api *restApi) Register(r *gin.RouterGroup) []web_utils.Controller {
	authController := api.NewAuthController()
	usersController := api.NewUsersController()
	var ctrls = []web_utils.Controller{
		authController.RegisterRoutes(r.Group("/auth")),
		usersController.RegisterRoutes(r.Group("/users")),
	}

	return ctrls
}

func (api *restApi) NewUsersController() *ctrl.UsersController {
	return &ctrl.UsersController{
		Users: api.server.App.Facades.Users,
		Http:  api.server.httpTool,
	}
}

func (api *restApi) NewAuthController() *ctrl.AuthController {
	usersFacade := api.server.App.Facades.Users
	return &ctrl.AuthController{
		Users:    usersFacade,
		Auth:     api.server.App.Facades.Auth,
		Http:     api.server.httpTool,
	}
}
