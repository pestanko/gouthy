package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/domain/apps"
	"github.com/pestanko/gouthy/app/web/api_errors"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

func NewAppController(tools *web_utils.HTTPTools) *AppController {
	return &AppController{
		Apps:  tools.App.Facades.Apps,
		Http:  tools,
	}
}

type AppController struct {
	Apps apps.Facade
	Http *web_utils.HTTPTools
}

func (ctrl *AppController) RegisterRoutes(r *gin.RouterGroup) {

	r.GET("", ctrl.List)
	r.POST("", ctrl.Create)
	r.GET("/:aid", ctrl.GetOne)
	r.PUT("/:aid", ctrl.Update)
	r.PATCH("/:aid", ctrl.Update)
	r.DELETE("/:aid", ctrl.Delete)
}

func (ctrl *AppController) List(context *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(context)
	list, err := ctrl.Apps.List(ctx, apps.ListParams{})
	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	ctrl.Http.JSON(ctx, 200, list)
}

func (ctrl *AppController) Create(context *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(context)

	var newApp apps.CreateDTO
	ginCtx := ctrl.Http.Gin(ctx)
	if err := ginCtx.Bind(&newApp); err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}
	app, err := ctrl.Apps.Create(ctx, &newApp)
	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	ctrl.Http.JSON(ctx, 201, app)
}

func (ctrl *AppController) GetOne(c *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(c)
	sid := ctrl.Http.Param(ctx, "aid")

	found, err := ctrl.Apps.GetByAnyId(ctx, sid)
	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	if found == nil {
		ctrl.Http.Fail(ctx, api_errors.NewAppNotFound().WithDetail(api_errors.ErrorDetail{
			"id": sid,
		}))
		return
	}

	ctrl.Http.JSON(ctx, 200, found)
}

func (ctrl *AppController) Update(c *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(c)
	sid := ctrl.Http.Param(ctx, "aid")

	found, err := ctrl.Apps.GetByAnyId(ctx, sid)

	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	if found == nil {
		ctrl.Http.Fail(ctx, api_errors.NewAppNotFound().WithDetail(api_errors.ErrorDetail{
			"id": sid,
		}))
		return
	}

	var updateApp apps.UpdateDTO
	if err := c.Bind(&updateApp); err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}
	user, err := ctrl.Apps.Update(ctx, found.ID, &updateApp)

	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	ctrl.Http.JSON(ctx, 201, user)
}

func (ctrl *AppController) Delete(context *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(context)
	sid := ctrl.Http.Param(ctx, "aid")

	found, err := ctrl.Apps.GetByAnyId(ctx, sid)
	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	if found == nil {
		ctrl.Http.Fail(ctx, api_errors.NewAppNotFound().WithDetail(api_errors.ErrorDetail{
			"id": sid,
		}))
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
