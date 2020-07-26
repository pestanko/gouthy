package auth

import (
	"github.com/pestanko/gouthy/app/domain/apps"
	"github.com/pestanko/gouthy/app/domain/jwtlib"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/shared"
	log "github.com/sirupsen/logrus"
)


type LoginIdentity struct {
	shared.LoggingIdentity
	UserId   string
	ClientId string
	Scopes   shared.Scopes
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

func CreateLoginIdentityForUserAndApp(user *users.UserDTO, app *apps.ApplicationDTO) *LoginIdentity {
	return NewLoginIdentity(user, app, []string{})
}

func NewLoginIdentity(user *users.UserDTO, app *apps.ApplicationDTO, scopes []string) *LoginIdentity {
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
