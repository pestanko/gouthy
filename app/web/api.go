package web

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
	controllers2 "github.com/pestanko/gouthy/app/web/controllers"
	"github.com/pestanko/gouthy/app/web/shared"
)

func RegisterApiControllers(app *core.GouthyApp, r *gin.RouterGroup) []shared.Controller {
	var controllers = []shared.Controller{
		controllers2.NewAuthController(app).RegisterRoutes(r.Group("/auth")),
		controllers2.NewUsersController(app).RegisterRoutes(r.Group("/users")),
		controllers2.NewMachinesController(app).RegisterRoutes(r.Group("/machines")),
		controllers2.NewEntitiesController(app).RegisterRoutes(r.Group("/entities")),
	}

	return controllers
}
