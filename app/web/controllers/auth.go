package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/domain/auth"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/web/api_errors"
	"github.com/pestanko/gouthy/app/web/web_utils"
	log "github.com/sirupsen/logrus"
)

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

func (ctrl *AuthController) RegisterRoutes(authRoute *gin.RouterGroup) web_utils.Controller {
	loginRoute := authRoute.Group("/login")
	authRoute.POST("/register", ctrl.Register)
	loginRoute.POST("/password", ctrl.LoginPassword)
	loginRoute.POST("/secret", ctrl.LoginSecret)

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
	if err := context.BindJSON(&loginDTO); err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}
	user, err := ctrl.Users.GetByUsername(ctx, loginDTO.Username)
	if err != nil {
		ctrl.Http.WriteErr(ctx, err)
		return
	}

	if user == nil {
		ctrl.Http.JSON(ctx, 401, gin.H{
			"status":   "not_found",
			"code":     401,
			"message":  "UserModel not found",
			"username": loginDTO.Username,
		})
		return
	}

	loginState := auth.NewLoginState(user.ID)
	loginState, err = ctrl.Auth.LoginUsernamePassword(ctx, loginState, loginDTO)

	if err != nil {
		ctrl.Http.Fail(ctx, api_errors.NewUnauthorizedError().WithError(err).WithDetail(api_errors.ErrorDetail{
			"username": loginDTO.Username,
		}))
		return
	}

	shared.GetLogger(ctx).WithFields(log.Fields{
		"loginState": loginState,
		"user_id": user.ID,
	}).Info("Creating login State")
	ctrl.Http.JSON(ctx, 200, loginState)
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
