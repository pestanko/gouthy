package auth

import (
	"github.com/pestanko/gouthy/app/apps"
	"github.com/pestanko/gouthy/app/jwtlib"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/users"
	log "github.com/sirupsen/logrus"
)

type LoginIdentity struct {
	UserId   string        `json:"user_id"`
	ClientId string        `json:"client_id"`
	Scopes   shared.Scopes `json:"scopes"`
}

func CreateLoginIdentityFromToken(token jwtlib.Jwt) *LoginIdentity {
	if token == nil {
		return nil
	}
	return &LoginIdentity{
		UserId:   token.UserId(),
		ClientId: token.AppId(), // same as an audience
		Scopes:   token.Scopes(),
	}
}

func CreateLoginIdentityForUserAndApp(user *users.UserDTO, app *apps.AppDTO) *LoginIdentity {
	return NewLoginIdentity(user, app, []string{})
}

func NewLoginIdentity(user *users.UserDTO, app *apps.AppDTO, scopes []string) *LoginIdentity {
	return &LoginIdentity{
		UserId:   user.ID.String(),
		ClientId: app.ClientId,
		Scopes:   scopes,
	}
}

func (id *LoginIdentity) LogFields() log.Fields {
	return log.Fields{
		"login_client_id": id.ClientId,
		"login_uid":       id.UserId,
		"login_scopes":    id.Scopes,
	}
}
