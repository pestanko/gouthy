package auth

import (
	"github.com/pestanko/gouthy/app/domain/apps"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/jwtlib"
	"github.com/pestanko/gouthy/app/shared"
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

func NewLoginIdentity(user *users.User, app *apps.Application, scopes []string) *LoginIdentity {
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
