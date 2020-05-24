package web_utils

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/infra"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/web/api_errors"
	uuid "github.com/satori/go.uuid"
)

type Controller interface {
	RegisterRoutes(router *gin.RouterGroup) Controller
}

type HTTPTools struct {
	App *infra.GouthyApp
}

func NewHTTPTools(app *infra.GouthyApp) *HTTPTools {
	return &HTTPTools{
		App: app,
	}
}

func (http *HTTPTools) NewControllerContext(gin *gin.Context) context.Context {
	identity := &UserIdentity{} // TODO

	ctx := shared.NewOperationContext()
	ctx = context.WithValue(ctx, "identity", identity)
	ctx = context.WithValue(ctx, "gin", gin)
	return ctx
}

type UserIdentity struct {
	UserId uuid.UUID
}

type ControllerContext struct {
	context.Context
	Gin *gin.Context
}

func (http *HTTPTools) JSON(ctx context.Context, code int, obj interface{}) {
	ginCtx := http.Gin(ctx)
	ginCtx.JSON(code, obj)
}

func (http *HTTPTools) Fail(ctx context.Context, err api_errors.ApiError) {
	http.JSON(ctx, err.Code(), err)
}

func (http *HTTPTools) WriteErr(ctx context.Context, err error) {
	http.Fail(ctx, api_errors.NewApiError().WithError(err))
}

func (http *HTTPTools) ParseUUID(ctx context.Context, id string) (uuid.UUID, error) {
	return uuid.FromString(id)
}

func (http *HTTPTools) Param(ctx context.Context, key string) string {
	ginCtx := http.Gin(ctx)
	return ginCtx.Param(key)
}

func (http *HTTPTools) Gin(ctx context.Context) *gin.Context {
	return ctx.Value("gin").(*gin.Context)
}
