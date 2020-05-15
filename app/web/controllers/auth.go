package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/domain/auth"
	"github.com/pestanko/gouthy/app/domain/entities"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/web/api_errors"
	"github.com/pestanko/gouthy/app/web/shared"
)

type AuthController struct {
	Entities entities.Facade
	Users    users.Facade
	Auth     auth.Facade
	Http     *shared.HTTPTools
}

func (ctrl *AuthController) RegisterRoutes(authRoute *gin.RouterGroup) shared.Controller {
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

func (ctrl *AuthController) LoginPassword(context *gin.Context) {
	ctx := ctrl.Http.NewControllerContext(context)
	var loginDTO auth.PasswordLoginDTO
	if err := ctx.Gin.BindJSON(&loginDTO); err != nil {
		ctx.WriteErr(err)
		return
	}
	user, err := ctrl.Users.GetByUsername(loginDTO.Username)
	if err != nil {
		ctx.WriteErr(err)
		return
	}

	if user != nil {
		ctx.Gin.JSON(401, gin.H{
			"status":   "not_found",
			"code":     401,
			"message":  "User not found",
			"username": loginDTO.Username,
		})
		return
	}
	tokens, err := ctrl.Auth.LoginUsernamePassword(loginDTO)

	if err != nil {
		ctx.Fail(api_errors.NewUnauthorizedError().WithError(err).WithDetail(api_errors.ErrorDetail{
			"username": loginDTO.Username,
		}))
	}
	ctx.JSON(200, tokens)
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
