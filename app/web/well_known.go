package web

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

type WellKnownController struct {
	Tools *web_utils.Tools
}

func NewWellKnownController(tools *web_utils.Tools) *WellKnownController {
	return &WellKnownController{Tools: tools}
}

func (c *WellKnownController) RegisterRoutes(route *gin.RouterGroup) {
	route.GET("/openid-configuration", c.OpenIdConfigurationEndpoint)
}

func (c *WellKnownController) OpenIdConfigurationEndpoint(context *gin.Context) {

}
