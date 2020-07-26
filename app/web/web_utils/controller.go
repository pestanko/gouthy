package web_utils

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/domain/apps"
	"github.com/pestanko/gouthy/app/domain/auth"
	"github.com/pestanko/gouthy/app/domain/jwtlib"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/infra"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/web/api_errors"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const CookieAccessToken = "JWT_ACCESS"
const CookieSessionToken = "JWT_SESSION"
const CookieRefreshToken = "JWT_REFRESH"

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
		UserId:   uuid.Nil.String(),
		ClientId: "default",
		Scopes:   []string{shared.ScopeUnauthorized, shared.ScopeAnonymous},
	}
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
	clientId := identity.ClientId

	if clientId == "" {
		clientId = apps.DefaultApplicationClientId
	}

	return tool.App.DI.Apps.Facade.GetByClientId(ctx, clientId)
}

func (tool *HTTPTools) GetIdentity(ctx context.Context) *auth.LoginIdentity {
	return ctx.Value("identity").(*auth.LoginIdentity)
}

func (tool *HTTPTools) ExtractJwt(ctx context.Context) (jwtlib.Jwt, error) {
	rawToken, err := extractJwkStringFromRequest(tool.Gin(ctx))
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"raw_token": rawToken, // TODO anonymize
		}).Error("")
		return nil, err
	}
	if rawToken == "" {
		shared.GetLogger(ctx).Debug("No token was found")
		return nil, err
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

func (tool *HTTPTools) GetLoggedInUser(ctx context.Context) *users.UserDTO {
	id := tool.GetIdentity(ctx)
	if id == nil || id.UserId == "" {
		return nil
	}
	dto, err := tool.App.Facades.Users.GetByAnyId(ctx, id.UserId)
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(id.LogFields()).Warn("Unable to get user")
	}
	return dto
}

func (tool *HTTPTools) extractIdentityFromRequest(ctx context.Context) *auth.LoginIdentity {
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

func extractJwkStringFromRequest(gin *gin.Context) (string, error) {
	// Extract from header
	authHeaderValue := extractAuthHeader(gin, "Bearer")
	if authHeaderValue != "" {
		return authHeaderValue, nil
	}

	// Extract from cookie
	sessionToken, err := gin.Cookie(CookieSessionToken)
	if err == nil {
		return sessionToken, nil
	}
	accessToken, err := gin.Cookie(CookieAccessToken)
	if err != http.ErrNoCookie {
		return "", err
	}
	return accessToken, err
}

func extractAuthHeader(gin *gin.Context, prefix string) string {
	authHeader := gin.GetHeader("Authorization")
	prefix = prefix + " "
	if authHeader != "" && strings.HasPrefix(authHeader, prefix) {
		return strings.TrimPrefix(authHeader, prefix)
	}
	return ""
}

func intoApiError(err error) api_errors.ApiError {
	switch v := err.(type) {
	case shared.GouthyError:
		return api_errors.FromGouthyError(v)
	default:
		return api_errors.NewApiError().WithError(err)
	}
}
