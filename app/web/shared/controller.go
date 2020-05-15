package shared

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/infra"
	"github.com/pestanko/gouthy/app/web/api_errors"
	uuid "github.com/satori/go.uuid"
)

type Controller interface {
	RegisterRoutes(router *gin.RouterGroup) Controller
}

type StandardResponses struct {

}

type HTTPTools struct {
	App *infra.GouthyApp
}

func NewHTTPTools(app *infra.GouthyApp) *HTTPTools {
	return &HTTPTools{
		App: app,
	}
}

func (http *HTTPTools) NewControllerContext(gin *gin.Context) *ControllerContext {
	return &ControllerContext{Gin: gin, App: http.App}
}


type ControllerContext struct {
	Gin       *gin.Context
	App       *infra.GouthyApp
	Responses StandardResponses
}

func (ctx *ControllerContext) JSON(code int, obj interface{}) {
	ctx.Gin.JSON(code, obj)
}

func (ctx *ControllerContext) Fail(err api_errors.ApiError) {
	ctx.JSON(err.Code(), err)
}

func (ctx *ControllerContext) WriteErr(err error) {
	ctx.Fail(api_errors.NewApiError().WithError(err))
}

func (ctx *ControllerContext) ParseUUID(id string) (uuid.UUID, error) {
	return uuid.FromString(id)
}

func (ctx *ControllerContext) Param(key string) string {
	return ctx.Gin.Param(key)
}


