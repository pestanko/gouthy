package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/web/api_errors"
	"github.com/pestanko/gouthy/app/web/shared"
)

type UsersController struct {
	Users users.Facade
	Http  *shared.HTTPTools
}

func (ctrl *UsersController) RegisterRoutes(r *gin.RouterGroup) shared.Controller {

	r.GET("", ctrl.List)
	r.POST("", ctrl.Create)
	r.GET("/:uid", ctrl.GetOne)
	r.PUT("/:uid", ctrl.Update)
	r.PATCH("/:uid", ctrl.Update)
	r.DELETE("/:uid", ctrl.Delete)
	r.POST("/:uid/password", ctrl.UpdatePassword)

	return ctrl
}

func (ctrl *UsersController) GetOne(context *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(context)
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
	ctx := ctrl.Http.NewControllerContext(context)
	listUsers, err := ctrl.Users.List()
	if err != nil {
		ctx.WriteErr(err)
		return
	}

	ctx.JSON(200, listUsers)
}

func (ctrl *UsersController) Delete(context *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(context)
	sid := ctx.Param("uid")

	foundUser, err := ctrl.Users.GetByAnyId(sid)
	if err != nil {
		ctx.WriteErr(err)
		return
	}

	err = ctrl.Users.Delete(foundUser.ID)
	ctx.Gin.Status(204)
}

func (ctrl *UsersController) Create(context *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(context)

	var newUser users.NewUserDTO
	if err := ctx.Gin.Bind(&newUser); err != nil {
		ctx.WriteErr(err)
		return
	}
	user, err := ctrl.Users.Create(&newUser)
	if err != nil {
		ctx.WriteErr(err)
		return
	}

	ctx.JSON(201, user)
}

func (ctrl *UsersController) Update(context *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(context)
	sid := ctx.Param("uid")

	foundUser, err := ctrl.Users.GetByAnyId(sid)

	if err != nil {
		ctx.WriteErr(err)
		return
	}

	var updateUser users.UpdateUserDTO
	if err := ctx.Gin.Bind(&updateUser); err != nil {
		ctx.WriteErr(err)
		return
	}
	user, err := ctrl.Users.Update(foundUser.ID, &updateUser)

	if err != nil {
		ctx.WriteErr(err)
		return
	}

	ctx.JSON(201, user)
	return
}

func (ctrl *UsersController) UpdatePassword(context *gin.Context) {

}

func (ctrl *UsersController) findUser(ctx *shared.ControllerContext, sid string) (*users.UserDTO, error) {
	user, err := ctrl.Users.GetByAnyId(sid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		ctx.Fail(api_errors.NewNotFound().WithMessage("User not found"))
	}
	return user, nil
}
