package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

func RegisterApiControllers(app *core.GouthyApp, r *gin.RouterGroup) []web_utils.Controller {
	var controllers = []web_utils.Controller{
		NewAuthController(app).RegisterRoutes(r.Group("/auth")),
		NewUsersController(app).RegisterRoutes(r.Group("/users")),
		NewMachinesController(app).RegisterRoutes(r.Group("/machines")),
		NewEntitiesController(app).RegisterRoutes(r.Group("/entities")),
	}

	return controllers
}
