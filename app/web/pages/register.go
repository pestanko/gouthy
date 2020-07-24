package pages

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/web/web_utils"
	"net/http"
)

func NewRegistrationController(tools *web_utils.HTTPTools) *RegistrationController {
	return &RegistrationController{
		Tools: tools,
	}
}

type RegistrationController struct {
	Tools *web_utils.HTTPTools
}

func (c *RegistrationController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/register", c.registrationPage)
	r.POST("/register", c.registrationPagePost)
}

func (c *RegistrationController) registrationPage(context *gin.Context) {
	ctx := c.Tools.NewControllerContext(context)
	c.Tools.HTML(ctx, http.StatusOK, "register.html", gin.H{})
}

func (c *RegistrationController) registrationPagePost(context *gin.Context) {
	ctx := c.Tools.NewControllerContext(context)
	c.Tools.HTML(ctx, http.StatusOK, "register.html", gin.H{})
}

