package auth

import (
	"context"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/infra/jwtlib"
	"github.com/pestanko/gouthy/app/shared"
	log "github.com/sirupsen/logrus"
)

type Facade interface {
	LoginUsernamePassword(ctx context.Context, loginState LoginState, pwd PasswordLoginDTO) (LoginState, error)
	LoginTOTP(ctx context.Context, state LoginState, totp TotpDTO) (LoginState, error)
	LoginUsingSecret(ctx context.Context, loginState LoginState, secret SecretLoginDTO) (LoginState, error)
}

type FacadeImpl struct {
	Users users.Repository
	Jwk   jwtlib.JwkRepository
}

func (auth *FacadeImpl) LoginTOTP(ctx context.Context, state LoginState, totp TotpDTO) (LoginState, error) {
	return nil, nil
}

func (auth *FacadeImpl) LoginUsernamePassword(ctx context.Context, loginState LoginState, pwd PasswordLoginDTO) (LoginState, error) {
	user, err := auth.Users.FindByUsername(ctx, pwd.Username)
	if err != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"username": pwd.Username,
		}).WithError(err).Debug("Unable to find user - error happened")
		return nil, err
	}
	flow := UcPasswordFlow{Users: auth.Users, user: user, password: pwd.Password}

	err = flow.Check()
	if err != nil {
		loginState.AddStep(NewLoginStep("UserPassword", Failed))
		shared.GetLogger(ctx).WithFields(log.Fields{
			"username": pwd.Username,
			"user_id": user.ID,
		}).Debug("User password check failed - error happened")
		return loginState, err
	}

	loginState.AddStep(NewLoginStep("UserPassword", Success))

	return loginState, nil
}

func (auth *FacadeImpl) LoginUsingSecret(ctx context.Context, loginState LoginState, secret SecretLoginDTO) (LoginState, error) {
	// Get Entity

	// check secret

	// create tokens
	return nil, nil
}

func NewAuthFacade(users users.Repository, inventory jwtlib.JwkRepository) Facade {
	return &FacadeImpl{Users: users, Jwk: inventory}
}

type PasswordLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SecretLoginDTO struct {
	Secret     string `json:"secret"`
	Codename   string `json:"codename"`
	EntityType string `json:"entity_type"`
}

type TotpDTO struct {
	TotpCode string
}
