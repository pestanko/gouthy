package web_utils

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
	uuid "github.com/satori/go.uuid"
)

type Controller interface {
	RegisterRoutes(router *gin.RouterGroup) Controller
}

type StandardResponses struct {

}

type HTTPTools struct {
	App *core.GouthyApp
}

func NewHTTPTools(app *core.GouthyApp) *HTTPTools {
	return &HTTPTools{
		App: app,
	}
}

func (http *HTTPTools) NewControllerContext(gin *gin.Context) *ControllerContext {
	return &ControllerContext{Gin: gin, App: http.App}
}


type ControllerContext struct {
	Gin *gin.Context
	App *core.GouthyApp
	Responses StandardResponses
}

func (ctx *ControllerContext) JSON(code int, obj interface{}) {
	ctx.Gin.JSON(code, obj)
}

func (ctx *ControllerContext) Fail(err ApiError) {
	ctx.JSON(err.Code, err)
}

func (ctx *ControllerContext) WriteError(err string, message string, code int) {
	ctx.Fail(NewApiError(err, message, code))
}

func (ctx *ControllerContext) WriteErr(err error) {
	ctx.Fail(NewApiError("server_error", err.Error(), 500))
}

func (ctx *ControllerContext) ParseUUID(id string) (uuid.UUID, error) {
	return uuid.FromString(id)
}

func (ctx *ControllerContext) Param(key string) string {
	return ctx.Gin.Param(key)
}


