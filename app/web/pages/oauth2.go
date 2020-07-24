package pages

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/web/web_utils"
	"net/http"
)

func NewOAuth2Controller(tools *web_utils.HTTPTools) *OAuth2Controller {
	return &OAuth2Controller{
		Tools: tools,
	}
}

type OAuth2Controller struct {
	Tools *web_utils.HTTPTools
}

func (c *OAuth2Controller) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/authorization-consent", c.authorizationConsentPage)

}

func (c *OAuth2Controller) authorizationConsentPage(context *gin.Context) {
	ctx := c.Tools.NewControllerContext(context)
	c.Tools.HTML(ctx, http.StatusOK, "authorization-consent.html", gin.H{})
}

