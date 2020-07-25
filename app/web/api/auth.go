package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/domain/auth"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/web/api_errors"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

func NewAuthController(tools *web_utils.HTTPTools) *AuthController {
	return &AuthController{
		Users: tools.App.DI.Users.Facade,
		Auth:  tools.App.DI.Auth.Facade,
		Http:  tools,
	}
}

type AuthController struct {
	Users users.Facade
	Auth  auth.Facade
	Http  *web_utils.HTTPTools
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func (c *AuthController) RegisterRoutes(r *gin.RouterGroup) {
	loginRoute := r.Group("/login")
	r.POST("/register", c.Register)
	loginRoute.POST("/password", c.LoginPassword)
	loginRoute.POST("/secret", c.LoginSecret)

	oauth2Route := r.Group("/oauth2")
	oauth2Route.GET("/authorize", c.OAuth2AuthorizeEndpoint)
	oauth2Route.POST("/token", c.OAuth2TokenEndpoint)
	oauth2Route.GET("/userinfo", c.OAuth2UserInfoEndpoint)
	oauth2Route.GET("/certs", c.OAuth2CertificatesEndpoint)
}

type PasswordLoginDTO struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (c *AuthController) LoginPassword(context *gin.Context) {
	ctx := c.Http.NewControllerContext(context)
	var credentials PasswordLoginDTO
	if err := context.BindJSON(&credentials); err != nil {
		c.Http.WriteErr(ctx, err)
		return
	}

	loginState, err := c.Auth.Login(ctx, auth.Credentials{
		Username: credentials.Username,
		Password: credentials.Password,
	})

	if err != nil {
		c.Http.Fail(ctx, api_errors.NewUnauthorizedError().WithError(err).WithDetail(api_errors.ErrorDetail{
			"username": credentials.Username,
		}))
		return
	}

	c.Http.JSON(ctx, 200, loginState)
}

func (c *AuthController) LoginSecret(context *gin.Context) {

}

func (c *AuthController) Register(context *gin.Context) {

}

func (c *AuthController) OAuth2AuthorizeEndpoint(context *gin.Context) {
	//ctx := c.Tools.NewControllerContext(context)
	//authorizationRequest := web_utils.OAuth2ParseAuthorizationRequest(context)

}

func (c *AuthController) OAuth2TokenEndpoint(context *gin.Context) {

}

func (c *AuthController) OAuth2UserInfoEndpoint(context *gin.Context) {

}

func (c *AuthController) OAuth2CertificatesEndpoint(context *gin.Context) {

}
