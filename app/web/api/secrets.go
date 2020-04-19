package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

type SecretsController struct {
	App  *core.GouthyApp
	http *web_utils.HTTPTools
}

func NewSecretsController(app *core.GouthyApp) *SecretsController {
	return &SecretsController{App: app, http: web_utils.NewHTTPTools(app)}
}

func (ctrl *SecretsController) RegisterRoutes(r *gin.RouterGroup) web_utils.Controller {
	r.GET("", ctrl.List)
	r.POST("", ctrl.Create)
	r.GET("/:sid", ctrl.GetOne)
	r.DELETE("/:sid", ctrl.Delete)
	r.PUT("/:sid", ctrl.Update)
	r.PATCH("/:sid", ctrl.Update)
	r.POST("/:sid/revoke", ctrl.Revoke)
	return ctrl
}

func (ctrl *SecretsController) List(context *gin.Context) {

}

func (ctrl *SecretsController) Create(context *gin.Context) {

}

func (ctrl *SecretsController) GetOne(context *gin.Context) {

}

func (ctrl *SecretsController) Delete(context *gin.Context) {

}

func (ctrl *SecretsController) Revoke(context *gin.Context) {

}

func (ctrl *SecretsController) Update(context *gin.Context) {

}
