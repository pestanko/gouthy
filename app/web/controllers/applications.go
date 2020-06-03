package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/domain/applications"
	"github.com/pestanko/gouthy/app/web/api_errors"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

type ApplicationsController struct {
	Apps applications.Facade
	Http *web_utils.HTTPTools
}

func (ctrl *ApplicationsController) RegisterRoutes(r *gin.RouterGroup) web_utils.Controller {

	r.GET("", ctrl.List)
	r.POST("", ctrl.Create)
	r.GET("/:aid", ctrl.GetOne)
	r.PUT("/:aid", ctrl.Update)
	r.PATCH("/:aid", ctrl.Update)
	r.DELETE("/:aid", ctrl.Delete)

	return ctrl
}

func (ctrl *ApplicationsController) List(context *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(context)
	list, err := ctrl.Apps.List(ctx, applications.ListParams{})
	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	ctrl.Http.JSON(ctx, 200, list)
}

func (ctrl *ApplicationsController) Create(context *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(context)

	var newApp applications.CreateDTO
	ginCtx := ctrl.Http.Gin(ctx)
	if err := ginCtx.Bind(&newApp); err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}
	user, err := ctrl.Apps.Create(ctx, &newApp)
	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	ctrl.Http.JSON(ctx, 201, user)
}

func (ctrl *ApplicationsController) GetOne(c *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(c)
	sid := ctrl.Http.Param(ctx, "aid")

	user, err := ctrl.findApp(ctx, sid)
	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	if user == nil {
		return
	}

	ctrl.Http.JSON(ctx, 200, user)
}

func (ctrl *ApplicationsController) Update(c *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(c)
	sid := ctrl.Http.Param(ctx, "aid")

	foundUser, err := ctrl.Apps.GetByAnyId(ctx, sid)

	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	var updateApp applications.UpdateDTO
	if err := c.Bind(&updateApp); err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}
	user, err := ctrl.Apps.Update(ctx, foundUser.ID, &updateApp)

	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	ctrl.Http.JSON(ctx, 201, user)
}

func (ctrl *ApplicationsController) Delete(context *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(context)
	sid := ctrl.Http.Param(ctx, "aid")

	found, err := ctrl.Apps.GetByAnyId(ctx, sid)
	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	err = ctrl.Apps.Delete(ctx, found.ID)
	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}
	ginCtx := ctrl.Http.Gin(ctx)
	ginCtx.Status(204)
}

func (ctrl *ApplicationsController) findApp(ctx context.Context, sid string) (*applications.ApplicationDTO, error) {
	found, err := ctrl.Apps.GetByAnyId(ctx, sid)
	if err != nil {
		return nil, err
	}
	if found == nil {
		ctrl.Http.Fail(ctx, api_errors.NewNotFound().WithMessage("Application not found"))
	}
	return found, nil
}
