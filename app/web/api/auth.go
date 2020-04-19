package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

type AuthController struct {
	App  *core.GouthyApp
	http *web_utils.HTTPTools
}

func NewAuthController(app *core.GouthyApp) *AuthController {
	return &AuthController{App: app, http: web_utils.NewHTTPTools(app)}
}

func (ctrl *AuthController) RegisterRoutes(authRoute *gin.RouterGroup) web_utils.Controller {
	loginRoute := authRoute.Group("/login")
	authRoute.POST("/register", ctrl.Register)
	loginRoute.POST("/login/password", ctrl.LoginPassword)
	loginRoute.POST("/login/secret", ctrl.LoginSecret)

	oauth2Route := authRoute.Group("/oauth2")
	oauth2Route.GET("/authorize", ctrl.OAuth2AuthorizeEndpoint)
	oauth2Route.POST("/token", ctrl.OAuth2TokenEndpoint)
	oauth2Route.GET("/userinfo", ctrl.OAuth2UserInfoEndpoint)
	oauth2Route.GET("/certs", ctrl.OAuth2CertificatesEndpoint)

	return ctrl
}

type PasswordLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ctrl *AuthController) LoginPassword(context *gin.Context) {
	ctx := ctrl.http.NewControllerContext(context)
	var loginDTO PasswordLoginDTO
	if err := ctx.Gin.BindJSON(&loginDTO); err != nil {
		ctx.WriteErr(err)
		return
	}
	username := loginDTO.Username
	user, err := ctrl.App.Services.Users.GetByUsername(username)
	if err != nil {
		ctx.WriteErr(err)
		return
	}

	if user != nil {
		ctx.Gin.JSON(401, gin.H{
			"status":   "not_found",
			"code":     401,
			"message":  "User not found",
			"username": username,
		})
		return
	}

	if user.CheckPassword(loginDTO.Password) {

	}

	response := ctx.App.Services.Auth.LoginByID(user.ID)
	ctx.JSON(200, response)
}

func (ctrl *AuthController) LoginSecret(context *gin.Context) {

}

func (ctrl *AuthController) Register(context *gin.Context) {

}

func (ctrl *AuthController) OAuth2AuthorizeEndpoint(context *gin.Context) {

}

func (ctrl *AuthController) OAuth2TokenEndpoint(context *gin.Context) {

}

func (ctrl *AuthController) OAuth2UserInfoEndpoint(context *gin.Context) {

}

func (ctrl *AuthController) OAuth2CertificatesEndpoint(context *gin.Context) {

}
