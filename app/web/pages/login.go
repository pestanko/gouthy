package pages

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/auth"
	"github.com/pestanko/gouthy/app/jwtlib"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/users"
	"github.com/pestanko/gouthy/app/web/api_errors"
	"github.com/pestanko/gouthy/app/web/web_utils"
	log "github.com/sirupsen/logrus"
)

func NewLoginPagesController(http *web_utils.Tools, users users.Facade, auth auth.Facade) *LoginPagesController {
	return &LoginPagesController{
		Http:  http,
		Users: users,
		Auth:  auth,
	}
}

type LoginPagesController struct {
	Http  *web_utils.Tools
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
	var cred loginCredentials

	if err := context.Bind(&cred); err != nil {
		c.Http.WriteErr(ctx, err)
		return
	}

	shared.GetLogger(ctx).WithFields(log.Fields{
		"username": cred.Username,
		"password": cred.Password, // TODO: Remove
		"rs":       cred.RedirectState,
		"state":    cred.State,
	}).Debug("Provided data")

	loginState, err := c.Auth.Login(ctx, auth.Credentials{
		Username: cred.Username,
		Password: cred.Password,
	})

	shared.GetLogger(ctx).WithFields(loginState.LogFields()).Info("Login handler result")

	if loginState != nil && loginState.IsOk() {
		c.doSuccessLogin(ctx, loginState, cred)
		return
	} else {
		c.doErrorLogin(ctx, cred, loginState, err)
	}
}

func (c *LoginPagesController) doSuccessLogin(ctx context.Context, state auth.LoginState, cred loginCredentials) {
	gctx := c.Http.Gin(ctx)
	user, err := c.Http.App.Facades.Users.GetByAnyId(ctx, state.UserID())
	if err != nil {
		shared.GetLogger(ctx).WithError(err).Info("Unable to get current user")
		c.Http.ErrorPage(ctx, web_utils.ErrorPageParams{
			Message: "Unable to get current user",
			Error:   err.Error(),
			Title:   "Error",
		})
		return
	}

	app, err := c.Http.GetCurrentApp(ctx)
	if err != nil {
		shared.GetLogger(ctx).WithError(err).Info("Unable to get current app")
		c.Http.ErrorPage(ctx, web_utils.ErrorPageParams{
			Message: "Unable to get current app",
			Error:   err.Error(),
			Title:   "Error",
		})
		return
	}

	identity := auth.NewLoginIdentity(user, app, []string{shared.ScopeUI, shared.ScopeSession})
	tokens, err := c.Auth.CreateSignedTokensFromLoginIdentity(ctx, identity)
	c.setTokensAsCookies(gctx, tokens)
	if err := c.Http.RedirectWithRedirectState(ctx, "/"); err != nil {
		shared.GetLogger(ctx).WithError(err).Info("Unable to redirect current app context")
		c.Http.ErrorPage(ctx, web_utils.ErrorPageParams{
			Message: "Unable to redirect",
			Error:   err.Error(),
			Title:   "Error",
		})
	}
}

func (c *LoginPagesController) setTokensAsCookies(gctx *gin.Context, tokens auth.SignedTokensDTO) {
	gctx.SetCookie(web_utils.CookieSessionToken, tokens.SessionToken, int(jwtlib.SessionTokenExpiration), "/", "", false, true)
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
