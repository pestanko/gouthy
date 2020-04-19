package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

type MachinesController struct {
	App  *core.GouthyApp
	http *web_utils.HTTPTools
}

func NewMachinesController(app *core.GouthyApp) *MachinesController {
	return &MachinesController{App: app, http: web_utils.NewHTTPTools(app)}
}

func (ctrl *MachinesController) RegisterRoutes(r *gin.RouterGroup) web_utils.Controller {
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
