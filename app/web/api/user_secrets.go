package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

type UserSecretsController struct {
	Users users.Facade
	Http  *web_utils.HTTPTools
}

func NewSecretsController(entitiesFacade users.Facade, http *web_utils.HTTPTools) *UserSecretsController {
	return &UserSecretsController{Users: entitiesFacade, Http: http}
}

func (ctrl *UserSecretsController) RegisterRoutes(r *gin.RouterGroup) web_utils.Controller {
	r.GET("", ctrl.List)
	r.POST("", ctrl.Create)
	r.GET("/:sid", ctrl.GetOne)
	r.DELETE("/:sid", ctrl.Delete)
	r.PUT("/:sid", ctrl.Update)
	r.PATCH("/:sid", ctrl.Update)
	r.POST("/:sid/revoke", ctrl.Revoke)
	return ctrl
}

func (ctrl *UserSecretsController) List(context *gin.Context) {

}

func (ctrl *UserSecretsController) Create(context *gin.Context) {

}

func (ctrl *UserSecretsController) GetOne(context *gin.Context) {

}

func (ctrl *UserSecretsController) Delete(context *gin.Context) {

}

func (ctrl *UserSecretsController) Revoke(context *gin.Context) {

}

func (ctrl *UserSecretsController) Update(context *gin.Context) {

}
