package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/app/services"
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
	r.POST("", c.Create)
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

func (c *UsersController) GetOne(ctx *web_utils.ControllerContext) error {
	sid := ctx.Param("id")

	user, err := c.findUser(ctx, sid)
	if err != nil {
		return err
	}

	if user == nil {
		return nil
	}

	ctx.JSON(200, user)

	return nil
}

func (c *UsersController) List(ctx *web_utils.ControllerContext) error {
	users, err := c.App.Services.Users.List()
	if err != nil {
		return err
	}

	ctx.JSON(200, users)

	return nil
}

func (c *UsersController) Update(ctx *web_utils.ControllerContext) error {
	sid := ctx.Param("id")

	foundUser, err := c.App.Services.Users.GetByAnyId(sid)

	if err != nil {
		return err
	}

	var updateUser services.UpdateUser
	if err := ctx.Gin.Bind(&updateUser); err != nil {
		return err
	}
	user, err := c.App.Services.Users.Update(foundUser.ID, &updateUser)
	if err != nil {
		return err
	}

	ctx.JSON(201, services.ConvertModelsToUserDTO(&user))
	return nil
}

func (c *UsersController) Delete(ctx *web_utils.ControllerContext) error {
	sid := ctx.Param("id")

	foundUser, err := c.App.Services.Users.GetByAnyId(sid)
	if err != nil {

		return err
	}

	err = c.App.Services.Users.Delete(foundUser.ID)
	ctx.Gin.Status(204)
	return err
}

func (c *UsersController) Create(ctx *web_utils.ControllerContext) error {
	var newUser services.NewUser
	if err := ctx.Gin.Bind(&newUser); err != nil {
		return err
	}
	user, err := c.App.Services.Users.Create(&newUser)
	if err != nil {
		return err
	}

	ctx.JSON(201, services.ConvertModelsToUserDTO(&user))
	return nil
}

func (c *UsersController) findUser(ctx *web_utils.ControllerContext, sid string) (*services.UserDTO, error) {
	user, err := c.App.Services.Users.GetByAnyId(sid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		ctx.WriteError("not_found", "User not found", 404)
	}
	return user, nil
}
