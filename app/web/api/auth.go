package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

type AuthController struct {
	App *core.GouthyApp
}

func CreateAuthController(app *core.GouthyApp) *AuthController {
	return &AuthController{App: app}
}

func (c *AuthController) RegisterRoutes(router *gin.RouterGroup) {
	route := c.newRouterController(router)

	route.POST("/login/password", c.LoginPassword)
	route.POST("/login/secret", c.LoginSecret)
	route.POST("/register", c.Register)
}

func (c *AuthController) newRouterController(router *gin.RouterGroup) web_utils.RouterWrapper {
	route := router.Group("/auth")

	r := web_utils.NewRouteController(route, c.App)
	return r
}

func (c *AuthController) LoginSecret(context *web_utils.ControllerContext) error {

	return nil
}

func (c *AuthController) Register(context *web_utils.ControllerContext) error {

	return nil
}

func (c *AuthController) LoginPassword(context *web_utils.ControllerContext) error {

	return nil
}
