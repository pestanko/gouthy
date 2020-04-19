package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/app/services"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

type UsersController struct {
	App  *core.GouthyApp
	http *web_utils.HTTPTools
}

func NewUsersController(app *core.GouthyApp) *UsersController {
	return &UsersController{App: app, http: web_utils.NewHTTPTools(app)}
}

func (ctrl *UsersController) RegisterRoutes(r *gin.RouterGroup) web_utils.Controller {

	r.GET("", ctrl.List)
	r.POST("", ctrl.Create)
	r.GET("/:uid", ctrl.GetOne)
	r.PUT("/:uid", ctrl.Update)
	r.PATCH("/:uid", ctrl.Update)
	r.DELETE("/:uid", ctrl.Delete)
	r.POST("/:uid/password", ctrl.UpdatePassword)

	NewSecretsController(ctrl.App).RegisterRoutes(r.Group("/:uid/secrets"))
	return ctrl
}

func (ctrl *UsersController) GetOne(context *gin.Context) {
	ctx := ctrl.http.NewControllerContext(context)
	sid := ctx.Param("id")

	user, err := ctrl.findUser(ctx, sid)
	if err != nil {
		ctx.WriteErr(err)
		return
	}

	if user == nil {
		return
	}

	ctx.JSON(200, user)
}

func (ctrl *UsersController) List(context *gin.Context) {
	ctx := ctrl.http.NewControllerContext(context)
	users, err := ctrl.App.Services.Users.List()
	if err != nil {
		ctx.WriteErr(err)
		return
	}

	ctx.JSON(200, users)
}

func (ctrl *UsersController) Delete(context *gin.Context) {
	ctx := ctrl.http.NewControllerContext(context)
	sid := ctx.Param("uid")

	foundUser, err := ctrl.App.Services.Users.GetByAnyId(sid)
	if err != nil {
		ctx.WriteErr(err)
		return
	}

	err = ctrl.App.Services.Users.Delete(foundUser.ID)
	ctx.Gin.Status(204)
}

func (ctrl *UsersController) Create(context *gin.Context) {
	ctx := ctrl.http.NewControllerContext(context)

	var newUser services.NewUser
	if err := ctx.Gin.Bind(&newUser); err != nil {
		ctx.WriteErr(err)
		return
	}
	user, err := ctrl.App.Services.Users.Create(&newUser)
	if err != nil {
		ctx.WriteErr(err)
		return
	}

	ctx.JSON(201, services.ConvertModelsToUserDTO(&user))
}

func (ctrl *UsersController) Update(context *gin.Context) {
	ctx := ctrl.http.NewControllerContext(context)
	sid := ctx.Param("uid")

	foundUser, err := ctrl.App.Services.Users.GetByAnyId(sid)

	if err != nil {
		ctx.WriteErr(err)
		return
	}

	var updateUser services.UpdateUser
	if err := ctx.Gin.Bind(&updateUser); err != nil {
		ctx.WriteErr(err)
		return
	}
	user, err := ctrl.App.Services.Users.Update(foundUser.ID, &updateUser)

	if err != nil {
		ctx.WriteErr(err)
		return
	}

	ctx.JSON(201, services.ConvertModelsToUserDTO(&user))
	return
}

func (ctrl *UsersController) UpdatePassword(context *gin.Context) {

}

func (ctrl *UsersController) findUser(ctx *web_utils.ControllerContext, sid string) (*services.UserDTO, error) {
	user, err := ctrl.App.Services.Users.GetByAnyId(sid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		ctx.WriteError("not_found", "User not found", 404)
	}
	return services.ConvertModelsToUserDTO(user), nil
}
