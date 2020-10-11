package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

func NewPasswordController(tools *web_utils.Tools) *PasswordController {
	return &PasswordController{
		Tools: tools,
	}
}

type PasswordController struct {
	Tools *web_utils.Tools
}

func (c *PasswordController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/forgot", c.ForgotPasswordPage)
	r.POST("/forgot", c.ForgotPasswordPagePost)
	r.GET("/forgot-code", c.ForgotPasswordCodePage)
	r.POST("/forgot-code", c.ForgotPasswordCodePagePost)
	r.GET("/change", c.ChangePasswordPage)
	r.POST("/change", c.ChangePasswordPagePost)
}

func (c *PasswordController) ForgotPasswordPage(context *gin.Context) {
	ctx := c.Tools.NewControllerContext(context)
	c.Tools.HTML(ctx, http.StatusOK, "forgot-password.html", gin.H{})
}

func (c *PasswordController) ForgotPasswordCodePage(context *gin.Context) {
	ctx := c.Tools.NewControllerContext(context)
	c.Tools.HTML(ctx, http.StatusOK, "forgot-password-code.html", gin.H{})
}

func (c *PasswordController) ChangePasswordPage(context *gin.Context) {
	ctx := c.Tools.NewControllerContext(context)
	c.Tools.HTML(ctx, http.StatusOK, "change-password.html", gin.H{})
}

func (c *PasswordController) ForgotPasswordPagePost(context *gin.Context) {

}

func (c *PasswordController) ForgotPasswordCodePagePost(context *gin.Context) {

}

func (c *PasswordController) ChangePasswordPagePost(context *gin.Context) {

}
