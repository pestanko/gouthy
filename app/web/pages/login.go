package pages

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/domain/auth"
	"github.com/pestanko/gouthy/app/domain/jwtlib"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/web/api_errors"
	"github.com/pestanko/gouthy/app/web/web_utils"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func NewLoginPagesController(http *web_utils.HTTPTools, users users.Facade, auth auth.Facade) *LoginPagesController {
	return &LoginPagesController{
		Http:  http,
		Users: users,
		Auth:  auth,
	}
}

type LoginPagesController struct {
	Http  *web_utils.HTTPTools
	Users users.Facade
	Auth  auth.Facade
}

func (c *LoginPagesController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/login", c.LoginPage)
	r.POST("/login", c.LoginPagePost)
}

func (c *LoginPagesController) LoginPage(context *gin.Context) {
	ctx := c.Http.NewControllerContext(context)
	c.Http.HTML(ctx, http.StatusOK, "login.html", gin.H{
		"state":          "my random state",
		"redirect_state": c.Http.GetRedirectState(ctx),
	})
}

type loginCredentials struct {
	Username      string `form:"username"`
	Password      string `form:"password"`
	RedirectState string `form:"_rs"`
	State         string `form:"state"`
}

func (c *LoginPagesController) LoginPagePost(context *gin.Context) {
	ctx := c.Http.NewControllerContext(context)
	var credentials loginCredentials

	if err := context.Bind(&credentials); err != nil {
		c.Http.WriteErr(ctx, err)
		return
	}

	loginState, err := c.Auth.Login(ctx, auth.Credentials{
		Username: credentials.Username,
		Password: credentials.Password,
	})

	if loginState != nil && loginState.IsOk() {
		c.doSuccessLogin(ctx, credentials, loginState)
		return
	} else {
		c.doErrorLogin(ctx, credentials, loginState, err)
	}
}

func (c *LoginPagesController) doSuccessLogin(ctx context.Context, credentials loginCredentials, state auth.LoginState) {
	gctx := c.Http.Gin(ctx)
	bytes, err := uuid.FromBytes([]byte(*state.UserID()))
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"username": credentials.Username,
		})
		printError(ctx, err, credentials)
		return
	}
	user, _ := c.Users.Get(ctx, bytes)
	app, err := c.Http.GetCurrentAppContext(ctx)
	tokens, err := c.Auth.CreateSignedTokensResponse(ctx, jwtlib.TokenCreateParams{
		User:   user,
		App:    app,
		Scopes: nil,
	})
	gctx.SetCookie()
}

func printError(ctx context.Context, err error, credentials loginCredentials) {

}

func (c *LoginPagesController) doErrorLogin(ctx context.Context, cred loginCredentials, state auth.LoginState, err error) {
	if err != nil {
		c.Http.Fail(ctx, api_errors.NewUnauthorizedError().WithError(err).WithDetail(api_errors.ErrorDetail{
			"username": cred.Username,
		}))
		return
	}

	c.Http.HTML(ctx, http.StatusOK, "login.html", gin.H{
		"redirect_state": c.Http.GetRedirectState(ctx),
		"state":          "my random state",
	})
}
