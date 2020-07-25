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
	Scopes   []string
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
		"client_id": id.ClientId,
		"uid":       id.UserId,
		"scopes":    id.Scopes,
	}
}
