package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
)

type UsersController struct {
	App *core.GouthyApp
}

func CreateUsersController(app *core.GouthyApp) *UsersController {
	return &UsersController{App: app}
}

func (c *UsersController) RegisterRoutes(router *gin.RouterGroup) {
	route := router.Group("/users")

	route.GET("", c.List)
	route.GET("/:id", c.GetOne)
	route.PUT("/:id", c.Update)
	route.PATCH("/:id", c.Update)
	route.DELETE("/:id", c.Delete)
}

func (c *UsersController) GetOne(context *gin.Context) {
	//id := context.Param("id")

}

func (c *UsersController) List(context *gin.Context) {

}

func (c *UsersController) Update(context *gin.Context) {

}

func (c *UsersController) Delete(context *gin.Context) {

}
