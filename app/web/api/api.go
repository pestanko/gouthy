package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/app/web"
)

func RegisterApiControllers(a *core.GouthyApp, r *gin.RouterGroup) []web.Controller {
	var controllers = []web.Controller{
		CreateAuthController(a),
		CreateUsersController(a),
	}

	for _, c := range controllers {
		c.RegisterRoutes(r)
	}

	return controllers
}
