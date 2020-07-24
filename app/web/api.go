package web

import (
	"github.com/pestanko/gouthy/app/web/api"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

func NewApiControllers(tools *web_utils.HTTPTools) *apisControllers {
	return &apisControllers{
		users: api.NewUsersController(tools),
		auth:  api.NewAuthController(tools),
		apps:  api.NewAppController(tools),
	}
}

type apisControllers struct {
	auth  *api.AuthController
	users *api.UsersController
	apps *api.AppController
}
