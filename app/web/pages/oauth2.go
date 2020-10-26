package pages

import (
	"github.com/pestanko/gouthy/app/auth"
	"github.com/pestanko/gouthy/app/shared"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

func NewOAuth2Controller(tools *web_utils.Tools) *OAuth2Controller {
	return &OAuth2Controller{
		Http: tools,
	}
}

type OAuth2Controller struct {
	Http *web_utils.Tools
}

func (c *OAuth2Controller) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/authorize", c.authorizationConsentPage)

}

func (c *OAuth2Controller) authorizationConsentPage(context *gin.Context) {
	ctx := c.Http.NewControllerContext(context)
	user := c.Http.GetLoggedInUser(ctx)

	if user == nil {
		c.Http.Redirect(ctx, "/pages/login?_rs=" + c.Http.EncodeCurrentUrl(ctx))
		return
	}

	request := auth.OAuth2AuthRequest{
		ClientId:      context.Query("client_id"),
		ResponseType:  context.Query("response_type"),
		RedirectUri:   context.Query("redirect_uri"),
		Scopes:        strings.Split(context.Query("scope"), " "),
		State:         context.Query("state"),
		Nonce:         context.Query("nonce"),
		PKCEChallenge: context.Query("code_challenge"),
		PKCEMethod:    context.Query("code_challenge_method"),
	}

	shared.GetLogger(ctx).WithField("req", request).Info("Parsed OAuth2 request")

	c.Http.HTML(ctx, http.StatusOK, "authorization-consent.html", gin.H{
		"user":      user,
		"client_id": request.ClientId,
	})
}
