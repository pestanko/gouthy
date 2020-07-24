package web_utils

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/domain/apps"
	"github.com/pestanko/gouthy/app/domain/jwtlib"
	"github.com/pestanko/gouthy/app/infra"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/web/api_errors"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strings"
)

type Controller interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type HTTPTools struct {
	App *infra.GouthyApp
}

func NewHTTPTools(app *infra.GouthyApp) *HTTPTools {
	return &HTTPTools{
		App: app,
	}
}

func (tool *HTTPTools) NewControllerContext(gin *gin.Context) context.Context {
	ctx := shared.NewOperationContext()
	identity := extractIdentityFromRequest(ctx, gin) // TODO

	ctx = context.WithValue(ctx, "identity", identity)
	ctx = context.WithValue(ctx, "gin", gin)
	return ctx
}

type UserIdentity struct {
	UserId uuid.UUID
	AppId uuid.UUID
}

type ControllerContext struct {
	context.Context
	Gin *gin.Context
}

func (tool *HTTPTools) JSON(ctx context.Context, code int, obj interface{}) {
	ginCtx := tool.Gin(ctx)
	ginCtx.JSON(code, obj)
}

func (tool *HTTPTools) HTML(ctx context.Context, code int, template string, params gin.H) {
	ginCtx := tool.Gin(ctx)
	ginCtx.HTML(code, template, params)
}

func (tool *HTTPTools) Fail(ctx context.Context, err api_errors.ApiError) {
	tool.JSON(ctx, err.Code(), err)
}

func (tool *HTTPTools) WriteErr(ctx context.Context, err error) {
	shared.LogError(shared.GetLogger(ctx), err)
	tool.Fail(ctx, intoApiError(err))
}

func (tool *HTTPTools) Param(ctx context.Context, key string) string {
	ginCtx := tool.Gin(ctx)
	return ginCtx.Param(key)
}

func (tool *HTTPTools) Gin(ctx context.Context) *gin.Context {
	return ctx.Value("gin").(*gin.Context)
}

func (tool *HTTPTools) GetRedirectState(ctx context.Context) string {
	const redirectState = "_rs"
	g := tool.Gin(ctx)
	postValue, ok := g.GetPostForm(redirectState)
	if ok {
		return postValue
	}

	return g.Query(redirectState)
}

func (tool *HTTPTools) EncodeCurrentUrl(ctx context.Context) string {
	g := tool.Gin(ctx)
	return EncodeUrlAndQuery(g.Request.URL)
}

func (tool *HTTPTools) Redirect(ctx context.Context, url string) {
	g := tool.Gin(ctx)
	g.Redirect(http.StatusFound, url)
}

func (tool *HTTPTools) RedirectWithRedirectState(ctx context.Context) error {
	state, err := DecodeRedirectState(tool.GetRedirectState(ctx))
	if err != nil {
		return err
	}
	tool.Redirect(ctx, state)
	return nil
}

func (tool *HTTPTools) GetCurrentAppContext(ctx context.Context) (*apps.ApplicationDTO, error) {
	identity := tool.GetIdentity(ctx)
	return tool.App.Facades.Apps.GetByAnyId(ctx, identity.AppId.String())
}

func (tool *HTTPTools) GetIdentity(ctx context.Context) UserIdentity {
	return ctx.Value("identity").(UserIdentity)
}



func (tool *HTTPTools) ExtractJwt(ctx context.Context) (*jwt.Token, error) {

}


func intoApiError(err error) api_errors.ApiError {
	switch v := err.(type) {
	case shared.GouthyError:
		return api_errors.FromGouthyError(v)
	default:
		return api_errors.NewApiError().WithError(err)
	}
}

func extractIdentityFromRequest(gin *gin.Context) UserIdentity {
	extractJwkFromRequest(gin)
}


func extractJwkFromRequest(gin *gin.Context) (jwtlib.Jwt, error) {
	// Extract from header
	authHeaderValue := extractAuthHeader(gin, "Bearer")
	if authHeaderValue != "" {
		return jwtlib.ParseString(authHeaderValue), nil
	}

	// Extract from cookie
}

func extractAuthHeader(gin *gin.Context, prefix string) string {
	authHeader := gin.GetHeader("Authorization")
	prefix = prefix + " "
	if authHeader != "" && strings.HasPrefix(authHeader, prefix) {
		return strings.TrimPrefix(authHeader, prefix)
	}
	return ""
}