package auth

import (
	"context"
	"github.com/pestanko/gouthy/app/domain/apps"
	"github.com/pestanko/gouthy/app/domain/users"
)

type CheckState struct {
	User        *users.User
	Application *apps.ApplicationDTO
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
		return loginState.AddStep(NewLoginStep(c.CheckName(), Error)), err
	}
	if !valid {
		return loginState.AddStep(NewLoginStep(c.CheckName(), Failed)), nil
	}
	return loginState.AddStep(NewLoginStep(c.CheckName(), Success)), nil
}

func (LoginCheckPassword) CheckName() string {
	return StepLoginPassword
}
