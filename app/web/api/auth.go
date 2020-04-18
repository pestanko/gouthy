package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
)

type AuthController struct {
	App *core.GouthyApp
}

func CreateAuthController(app *core.GouthyApp) *AuthController {
	return &AuthController{App: app}
}

func (c *AuthController) RegisterRoutes(router *gin.RouterGroup) {
	route := router.Group("/auth")

	route.POST("/login/password", c.LoginPassword)
	route.POST("/login/secret", c.LoginSecret)
	route.POST("/register", c.Register)
}

func (c *AuthController) LoginSecret(context *gin.Context) {

}

func (c *AuthController) Register(context *gin.Context) {

}

func (c *AuthController) LoginPassword(context *gin.Context) {

}
