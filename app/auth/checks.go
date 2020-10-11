package auth

import (
	"context"
	"github.com/pestanko/gouthy/app/apps"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/users"
)

type CheckState struct {
	User        *users.User
	Application *apps.AppDTO
	Password    string
	Secret      string
}

type LoginCheck interface {
	Check(ctx context.Context, loginState LoginState, checkState CheckState) (LoginState, error)
	CheckName() string
}

//
// loginPage Check password
//
type LoginCheckPassword struct {
	passwordService users.PasswordService
}

func NewLoginCheckPassword(service users.PasswordService) LoginCheck {
	return &LoginCheckPassword{passwordService: service}
}

func (c *LoginCheckPassword) Check(ctx context.Context, loginState LoginState, checkState CheckState) (LoginState, error) {
	valid, err := c.passwordService.CheckPassword(ctx, checkState.User, checkState.Password)
	if err != nil {
		shared.GetLogger(ctx).WithFields(loginState.LogFields()).WithError(err).Warning("Pwd check - error")
		return loginState.AddStep(NewLoginStep(c.CheckName(), Error)), err
	}
	if !valid {
		shared.GetLogger(ctx).WithFields(loginState.LogFields()).Debug("Pwd check - fail")
		return loginState.AddStep(NewLoginStep(c.CheckName(), Failed)), nil
	}
	shared.GetLogger(ctx).WithFields(loginState.LogFields()).Debug("Pwd check - success")

	return loginState.AddStep(NewLoginStep(c.CheckName(), Success)), nil
}

func (LoginCheckPassword) CheckName() string {
	return StepLoginPassword
}
