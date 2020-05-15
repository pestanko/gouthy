package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/domain/entities"
	"github.com/pestanko/gouthy/app/web/shared"
)

type EntitiesController struct {
	Entities entities.Facade
	Http     *shared.HTTPTools
}

func (ctrl *EntitiesController) RegisterRoutes(r *gin.RouterGroup) shared.Controller {
	r.GET("", ctrl.List)
	r.POST("", ctrl.Create)
	r.GET("/:id", ctrl.GetOne)
	r.PUT("/:id", ctrl.Update)
	r.PATCH("/:id", ctrl.Update)
	r.DELETE("/:id", ctrl.Delete)

	secrets := NewSecretsController(ctrl.Entities, ctrl.Http)
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
