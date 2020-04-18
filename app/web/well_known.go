package web

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

type WellKnownController struct {
	App *core.GouthyApp
}

func CreateWellKnownController(app *core.GouthyApp) *WellKnownController {
	return &WellKnownController{App: app}
}

func RegisterWellKnownControllers(a *core.GouthyApp, r *gin.RouterGroup) []web_utils.Controller {
	var controllers = []web_utils.Controller{
		CreateWellKnownController(a),
	}

	for _, c := range controllers {
		c.RegisterRoutes(r)
	}

	return controllers
}

func (c *WellKnownController) RegisterRoutes(route *gin.RouterGroup) {
	route.GET("/openid-configuration", c.OpenIdConfigurationEndpoint)
}

func (c *WellKnownController) OpenIdConfigurationEndpoint(context *gin.Context) {

}
