package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/web/api_errors"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

func NewUsersController(tools *web_utils.HTTPTools) *UsersController {
	return &UsersController{
		Users: tools.App.DI.Users.Facade,
		Http:  tools,
	}
}

type UsersController struct {
	Users users.Facade
	Http  *web_utils.HTTPTools
}

func (ctrl *UsersController) RegisterRoutes(r *gin.RouterGroup) {

	r.GET("", ctrl.List)
	r.POST("", ctrl.Create)
	r.GET("/:uid", ctrl.GetOne)
	r.PUT("/:uid", ctrl.Update)
	r.PATCH("/:uid", ctrl.Update)
	r.DELETE("/:uid", ctrl.Delete)
	r.POST("/:uid/password", ctrl.UpdatePassword)
}

func (ctrl *UsersController) GetOne(context *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(context)
	sid := ctrl.Http.Param(ctx, "uid")

	user, err := ctrl.Users.GetByAnyId(ctx, sid)
	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	if user == nil {
		ctrl.Http.Fail(ctx, api_errors.NewUserNotFound().WithDetail(api_errors.ErrorDetail{
			"id": sid,
		}))
	} else {
		ctrl.Http.JSON(ctx, 200, user)

	}

}

func (ctrl *UsersController) List(context *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(context)
	list, err := ctrl.Users.List(ctx, users.ListParams{})
	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	ctrl.Http.JSON(ctx, 200, list)
}

func (ctrl *UsersController) Delete(context *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(context)
	sid := ctrl.Http.Param(ctx, "uid")

	found, err := ctrl.Users.GetByAnyId(ctx, sid)
	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	if found == nil {
		ctrl.Http.Fail(ctx, api_errors.NewUserNotFound().WithDetail(api_errors.ErrorDetail{
			"id": sid,
		}))
		return
	}

	err = ctrl.Users.Delete(ctx, found.ID)
	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}
	ginCtx := ctrl.Http.Gin(ctx)
	ginCtx.Status(204)
}

func (ctrl *UsersController) Create(context *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(context)

	var newUser users.CreateDTO
	ginCtx := ctrl.Http.Gin(ctx)
	if err := ginCtx.Bind(&newUser); err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}
	user, err := ctrl.Users.Create(ctx, &newUser)
	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	ctrl.Http.JSON(ctx, 201, user)
}

func (ctrl *UsersController) Update(c *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(c)
	sid := ctrl.Http.Param(ctx, "uid")

	found, err := ctrl.Users.GetByAnyId(ctx, sid)

	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	if found == nil {
		ctrl.Http.Fail(ctx, api_errors.NewUserNotFound().WithDetail(api_errors.ErrorDetail{
			"id": sid,
		}))
		return
	}

	var updateUser users.UpdateDTO
	if err := c.Bind(&updateUser); err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}
	user, err := ctrl.Users.Update(ctx, found.ID, &updateUser)

	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	ctrl.Http.JSON(ctx, 201, user)
}

func (ctrl *UsersController) UpdatePassword(c *gin.Context) {

}
