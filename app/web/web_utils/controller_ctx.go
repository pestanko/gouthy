package web_utils

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
	uuid "github.com/satori/go.uuid"
)

type ControllerContext struct {
	Gin *gin.Context
	App *core.GouthyApp
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
