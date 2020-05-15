package web

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/infra"
	"github.com/pestanko/gouthy/app/web/shared"
)

type WellKnownController struct {
	App *infra.GouthyApp
}

func CreateWellKnownController(app *infra.GouthyApp) *WellKnownController {
	return &WellKnownController{App: app}
}

func RegisterWellKnownControllers(a *infra.GouthyApp, r *gin.RouterGroup) []shared.Controller {
	var controllers = []shared.Controller{
		CreateWellKnownController(a),
	}

	for _, c := range controllers {
		c.RegisterRoutes(r)
	}

	return controllers
}

func (c *WellKnownController) RegisterRoutes(route *gin.RouterGroup) shared.Controller {
	route.GET("/openid-configuration", c.OpenIdConfigurationEndpoint)
	return c
}

func (c *WellKnownController) OpenIdConfigurationEndpoint(context *gin.Context) {

}
