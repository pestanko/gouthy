package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/auth"
	"github.com/pestanko/gouthy/app/domain/apps"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

func NewInternalController(tools *web_utils.Tools) *InternalController {
	return &InternalController{
		Users: tools.App.DI.Users.Facade,
		Auth:  tools.App.DI.Auth.Facade,
		Http:  tools,
	}
}

type InternalController struct {
	Users users.Facade
	Auth  auth.Facade
	Http  *web_utils.Tools
}

func (c *InternalController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/debug", c.debugInfo)
}

func (c *InternalController) debugInfo(context *gin.Context) {
	ctx := c.Http.NewControllerContext(context)

	debugData, err := c.getDebugData(ctx)
	if err != nil {
		c.Http.WriteErr(ctx, err)
		return
	}

	c.Http.JSON(ctx, http.StatusOK, debugData)
}

type debugData struct {
	User           *users.UserDTO         `json:"User"`
	App            *apps.AppDTO           `json:"app"`
	Identity       auth.LoginIdentity     `json:"identity"`
	TokenClaims    map[string]interface{} `json:"jwt_claims"`
	RawToken       string                 `json:"jwt_raw"`
	TokenHeader    map[string]interface{} `json:"jwt_headers"`
	RequestHeaders map[string][]string    `json:"request_headers"`
}

func (c *InternalController) getDebugData(ctx context.Context) (debugData, error) {
	id := c.Http.GetIdentity(ctx)
	token, _ := c.Http.ExtractJwt(ctx)
	user, app, _ := c.Http.GetCurrentUserAndApp(ctx)

	data := debugData{
		User:           user,
		App:            app,
		Identity:       *id,
		RawToken:       "Not logged in",
		RequestHeaders: c.Http.Gin(ctx).Request.Header,
	}

	if token != nil {
		data.RawToken = token.Raw()
		data.TokenHeader = token.RawHeader()
		data.TokenClaims = token.Claims()
	}
	return data, nil
}
