package web

import (
	"github.com/gin-gonic/gin"
	ctrl "github.com/pestanko/gouthy/app/web/controllers"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

type RestApi struct {
	server *WebServer
}

func NewRestApi(server *WebServer) *RestApi {
	return &RestApi{server: server}
}

func (api *RestApi) Register(r *gin.RouterGroup) []web_utils.Controller {
	authController := api.NewAuthController()
	usersController := api.NewUsersController()
	var ctrls = []web_utils.Controller{
		authController.RegisterRoutes(r.Group("/auth")),
		usersController.RegisterRoutes(r.Group("/users")),
	}

	return ctrls
}

func (api *RestApi) NewUsersController() *ctrl.UsersController {
	return &ctrl.UsersController{
		Users: api.server.App.Facades.Users,
		Http:  api.server.httpTool,
	}
}

func (api *RestApi) NewAuthController() *ctrl.AuthController {
	usersFacade := api.server.App.Facades.Users
	return &ctrl.AuthController{
		Users:    usersFacade,
		Auth:     api.server.App.Facades.Auth,
		Http:     api.server.httpTool,
	}
}
