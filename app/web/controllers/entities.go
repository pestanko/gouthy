package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/app/web/shared"
)

type EntitiesController struct {
	App  *core.GouthyApp
	http *shared.HTTPTools
}

func NewEntitiesController(app *core.GouthyApp) *EntitiesController {
	return &EntitiesController{App: app, http: shared.NewHTTPTools(app)}
}

func (ctrl *EntitiesController) RegisterRoutes(r *gin.RouterGroup) shared.Controller {
	r.GET("", ctrl.List)
	r.POST("", ctrl.Create)
	r.GET("/:id", ctrl.GetOne)
	r.PUT("/:id", ctrl.Update)
	r.PATCH("/:id", ctrl.Update)
	r.DELETE("/:id", ctrl.Delete)

	secrets := NewSecretsController(ctrl.App)
	secrets.RegisterRoutes(r.Group("/:id/secrets"))
	return ctrl
}

func (ctrl *EntitiesController) List(context *gin.Context) {

}

func (ctrl *EntitiesController) Create(context *gin.Context) {

}

func (ctrl *EntitiesController) GetOne(context *gin.Context) {

}

func (ctrl *EntitiesController) Update(context *gin.Context) {

}

func (ctrl *EntitiesController) Delete(context *gin.Context) {

}
