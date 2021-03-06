package web

import (
	"github.com/pestanko/gouthy/app/web/pages"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

func NewPagesControllers(tools *web_utils.Tools) *pagesControllers {
	return &pagesControllers{
		login: pages.NewLoginPagesController(tools,
			tools.App.DI.Users.Facade,
			tools.App.DI.Auth.Facade),
		password:     pages.NewPasswordController(tools),
		registration: pages.NewRegistrationController(tools),
		oauth2:       pages.NewOAuth2Controller(tools),
		debug:        pages.NewDebugController(tools),
		index:        pages.NewIndexController(tools),
	}
}

type pagesControllers struct {
	login        *pages.LoginPagesController
	password     *pages.PasswordController
	registration *pages.RegistrationController
	oauth2       *pages.OAuth2Controller
	debug        *pages.DebugController
	index        *pages.IndexController
}
