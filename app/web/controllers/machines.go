package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/domain/machines"
	"github.com/pestanko/gouthy/app/web/shared"
)

type MachinesController struct {
	Machines    *machines.Facade
	Http        *shared.HTTPTools
}

func (ctrl *MachinesController) RegisterRoutes(r *gin.RouterGroup) shared.Controller {
	r.GET("", ctrl.List)
	r.POST("", ctrl.Create)
	r.GET("/:id", ctrl.GetOne)
	r.PUT("/:id", ctrl.Update)
	r.PATCH("/:id", ctrl.Update)
	r.DELETE("/:id", ctrl.Delete)

	return ctrl
}

func (ctrl *MachinesController) List(context *gin.Context) {

}

func (ctrl *MachinesController) Create(context *gin.Context) {

}

func (ctrl *MachinesController) GetOne(context *gin.Context) {

}

func (ctrl *MachinesController) Update(context *gin.Context) {

}

func (ctrl *MachinesController) Delete(context *gin.Context) {

}
