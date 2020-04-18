package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

type UsersController struct {
	App *core.GouthyApp
}

func CreateUsersController(app *core.GouthyApp) *UsersController {
	return &UsersController{App: app}
}

func (c *UsersController) RegisterRoutes(router *gin.RouterGroup) {
	r := c.newRouterController(router)

	r.GET("", c.List)
	r.GET("/:id", c.GetOne)
	r.PUT("/:id", c.Update)
	r.PATCH("/:id", c.Update)
	r.DELETE("/:id", c.Delete)
}

func (c *UsersController) newRouterController(router *gin.RouterGroup) web_utils.RouterWrapper {
	route := router.Group("/users")

	r := web_utils.NewRouteController(route, c.App)
	return r
}

func (c *UsersController) GetOne(context *web_utils.ControllerContext) error {
	//id := context.Param("id")

	return nil
}

func (c *UsersController) List(context *web_utils.ControllerContext) error {

	users, err := c.App.Services.Users.List()
	if err != nil {
		context.WriteErr(err)
	}

	context.JSON(200, users)

	return nil
}

func (c *UsersController) Update(context *web_utils.ControllerContext) error {

	return nil

}

func (c *UsersController) Delete(context *web_utils.ControllerContext) error {

	return nil
}
