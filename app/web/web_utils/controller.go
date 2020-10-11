package web_utils

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/apps"
	"github.com/pestanko/gouthy/app/auth"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/app/jwtlib"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/users"
	"github.com/pestanko/gouthy/app/web/api_errors"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const CookieAccessToken = "JWT_ACCESS"
const CookieSessionToken = "JWT_SESSION"
const CookieRefreshToken = "JWT_REFRESH"

type Controller interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type Tools struct {
	App *core.GouthyApp
}

func NewHTTPTools(app *core.GouthyApp) *Tools {
	return &Tools{
		App: app,
	}
}

func (tool *Tools) NewControllerContext(gin *gin.Context) context.Context {
	ctx := shared.NewContextWithConfiguration(tool.App.Config)
	ctx = context.WithValue(ctx, "gin", gin)

	identity := tool.extractIdentityFromRequest(ctx) // TODO
	if identity == nil {
		identity = defaultAnonymousIdentity(ctx)
	}
	ctx = context.WithValue(ctx, "identity", identity)
	shared.GetLogger(ctx).Debug("Created a new context with identity")
	return ctx
}

func defaultAnonymousIdentity(ctx context.Context) *auth.LoginIdentity {
	return &auth.LoginIdentity{
		UserId:   "",
		ClientId: "default",
		Scopes:   []string{shared.ScopeUnauthorized, shared.ScopeAnonymous},
	}
}

type ControllerContext struct {
	context.Context
	Gin *gin.Context
}

func (tool *Tools) JSON(ctx context.Context, code int, obj interface{}) {
	ginCtx := tool.Gin(ctx)
	ginCtx.JSON(code, obj)
}

func (tool *Tools) HTML(ctx context.Context, code int, template string, params gin.H) {
	ginCtx := tool.Gin(ctx)
	ginCtx.HTML(code, template, params)
}

func (tool *Tools) Fail(ctx context.Context, err api_errors.ApiError) {
	tool.JSON(ctx, err.Code(), err)
}

func (tool *Tools) WriteErr(ctx context.Context, err error) {
	shared.LogError(shared.GetLogger(ctx), err)
	tool.Fail(ctx, intoApiError(err))
}

func (tool *Tools) Param(ctx context.Context, key string) string {
	ginCtx := tool.Gin(ctx)
	return ginCtx.Param(key)
}

func (tool *Tools) Gin(ctx context.Context) *gin.Context {
	return ctx.Value("gin").(*gin.Context)
}

func (tool *Tools) GetRedirectState(ctx context.Context) string {
	const redirectState = "_rs"
	g := tool.Gin(ctx)
	postValue, ok := g.GetPostForm(redirectState)
	if ok {
		return postValue
	}

	return g.Query(redirectState)
}

func (tool *Tools) EncodeCurrentUrl(ctx context.Context) string {
	g := tool.Gin(ctx)
	return EncodeUrlAndQuery(g.Request.URL)
}

func (tool *Tools) Redirect(ctx context.Context, url string) {
	g := tool.Gin(ctx)
	g.Redirect(http.StatusFound, url)
}

func (tool *Tools) RedirectWithRedirectState(ctx context.Context, defaultRedirect string) error {
	redirect, err := DecodeRedirectState(tool.GetRedirectState(ctx))
	if err != nil {
		return err
	}
	if redirect == "" {
		redirect = defaultRedirect
	}
	tool.Redirect(ctx, redirect)
	return nil
}

func (tool *Tools) GetCurrentApp(ctx context.Context) (*apps.AppDTO, error) {
	identity := tool.GetIdentity(ctx)
	clientId := identity.ClientId

	if clientId == "" {
		clientId = apps.DefaultApplicationClientId
	}

	return tool.App.DI.Apps.Facade.GetByClientId(ctx, clientId)
}

func (tool *Tools) GetIdentity(ctx context.Context) *auth.LoginIdentity {
	return ctx.Value("identity").(*auth.LoginIdentity)
}

func (tool *Tools) ExtractJwt(ctx context.Context) (jwtlib.Jwt, error) {
	rawToken := extractJwkStringFromRequest(ctx, tool.Gin(ctx))

	if rawToken == "" {
		shared.GetLogger(ctx).Debug("No token was found")
		return nil, nil
	}
	token, err := tool.App.Facades.Auth.ParseAndValidateJwt(ctx, rawToken)
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"raw_token": rawToken, // TODO anonymize
		}).Error("Unable to parse or verify the token")
		return token, err
	}
	return token, nil
}

func (tool *Tools) GetLoggedInUser(ctx context.Context) *users.UserDTO {
	id := tool.GetIdentity(ctx)
	if id == nil || id.UserId == "" || uuid.FromStringOrNil(id.UserId) == uuid.Nil {
		return nil
	}

	dto, err := tool.App.Facades.Users.GetByAnyId(ctx, id.UserId)
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(id.LogFields()).Warn("Unable to get user")
		return nil
	}
	return dto
}

func (tool *Tools) extractIdentityFromRequest(ctx context.Context) *auth.LoginIdentity {
	token, _ := tool.ExtractJwt(ctx)
	if token != nil {
		identity, err := tool.App.Facades.Auth.CreateLoginIdentityFromToken(ctx, token)
		if err != nil {
			shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
				"jti":       token.ID(),
				"uid":       token.UserId(),
				"client_id": token.ClientId(),
			}).Error("Unable to create identity")
		}
		return identity
	}

	return nil
}

type ErrorPageParams struct {
	Message string
	Error   string
	Title   string
}

func (tool *Tools) ErrorPage(ctx context.Context, params ErrorPageParams) {
	tool.HTML(ctx, http.StatusInternalServerError, "error.html", gin.H{
		"message": params.Message,
		"title":   params.Title,
		"error":   params.Error,
	})
}

func (tool *Tools) GetCurrentUserAndApp(ctx context.Context) (user *users.UserDTO, app *apps.AppDTO, err error) {
	id := tool.GetIdentity(ctx)
	if id == nil {
		return user, app, fmt.Errorf("no identity was found")
	}

	if id.UserId != "" && uuid.FromStringOrNil(id.UserId) != uuid.Nil {
		user, err = tool.App.Facades.Users.GetByAnyId(ctx, id.UserId)
		if err != nil {
			return
		}
	}

	app, err = tool.App.Facades.Apps.GetByClientId(ctx, id.ClientId)

	return
}

func extractJwkStringFromRequest(ctx context.Context, gin *gin.Context) string {
	// Extract from header
	authHeaderValue := extractAuthHeader(gin, "Bearer")
	if authHeaderValue != "" {
		return authHeaderValue
	}

	cookieNames := []string{CookieSessionToken, CookieAccessToken}

	for _, name := range cookieNames {
		value := extractCookie(ctx, gin, name)
		if name != "" {
			return value
		}
	}
	return ""
}

func extractAuthHeader(gin *gin.Context, prefix string) string {
	authHeader := gin.GetHeader("Authorization")
	prefix = prefix + " "
	if authHeader != "" && strings.HasPrefix(authHeader, prefix) {
		return strings.TrimPrefix(authHeader, prefix)
	}
	return ""
}

func extractCookie(ctx context.Context, g *gin.Context, name string) string {
	logEntry := shared.GetLogger(ctx).WithFields(log.Fields{
		"cookie_name": name,
	})
	value, err := g.Cookie(name)
	if err == http.ErrNoCookie {
		logEntry.Debug("No cookie found with specified name")
		return ""
	}
	if err != nil {
		logEntry.WithError(err).Warning("Failed to extract cookie")
		return ""
	}

	logEntry.Debug("Cookie found with specified name")
	return value
}

func intoApiError(err error) api_errors.ApiError {
	switch v := err.(type) {
	case shared.AppError:
		return api_errors.FromGouthyError(v)
	default:
		return api_errors.NewApiError().WithError(err)
	}
}
