package auth

import (
	"context"
	"github.com/pestanko/gouthy/app/domain/apps"
	"github.com/pestanko/gouthy/app/domain/jwtlib"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/shared"
	log "github.com/sirupsen/logrus"
)

type SignedTokensDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func (d *SignedTokensDTO) Serialize() string {
	return shared.ToJSONIndent(d)
}

type Facade interface {
	LoginUsernamePassword(ctx context.Context, loginState LoginState, pwd PasswordLoginDTO) (LoginState, error)
	LoginTOTP(ctx context.Context, state LoginState, totp TotpDTO) (LoginState, error)
	LoginUsingSecret(ctx context.Context, loginState LoginState, secret SecretLoginDTO) (LoginState, error)

	CreateSignedTokensResponse(ctx context.Context, params jwtlib.TokenCreateParams) (SignedTokensDTO, error)
}

type facadeImpl struct {
	users           users.Repository
	JwkService      jwtlib.JwkService
	JwtService      jwtlib.JwtService
	passwordService users.PasswordService
}

func (auth *facadeImpl) GenerateNewJwk(ctx context.Context) error {
	return auth.JwkService.GenerateNew(ctx)
}

func (auth *facadeImpl) ListJwks(ctx context.Context) ([]jwtlib.Jwk, error) {
	return auth.JwkService.List(ctx)
}

func (auth *facadeImpl) CreateSignedTokensResponse(ctx context.Context, params jwtlib.TokenCreateParams) (SignedTokensDTO, error) {
	result := SignedTokensDTO{
		ExpiresIn: "3600",
		TokenType: "Bearer",
	}

	access, err := auth.JwtService.CreateSignedAccessToken(ctx, params)
	if err != nil {
		return SignedTokensDTO{}, err
	}

	if access != nil {
		result.AccessToken = access.Signature
	}

	refresh, err := auth.JwtService.CreateSignedRefreshToken(ctx, params)
	if err != nil {
		return SignedTokensDTO{}, err
	}

	if refresh != nil {
		result.RefreshToken = refresh.Signature
	}

	id, err := auth.JwtService.CreateSignedIdToken(ctx, params)
	if err != nil {
		return SignedTokensDTO{}, err
	}

	if id != nil {
		result.IdToken = id.Signature
	}
	return result, nil
}

func (auth *facadeImpl) LoginTOTP(ctx context.Context, state LoginState, totp TotpDTO) (LoginState, error) {
	return nil, nil
}

func (auth *facadeImpl) LoginUsernamePassword(ctx context.Context, loginState LoginState, pwd PasswordLoginDTO) (LoginState, error) {
	user, err := auth.users.QueryOne(ctx, users.FindQuery{Username: pwd.Username})

	if err != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"username": pwd.Username,
		}).WithError(err).Debug("Unable to find user - error happened")
		return nil, err
	}

	check := NewLoginCheckPassword(auth.passwordService)

	entry := shared.GetLogger(ctx).WithFields(log.Fields{
		"username": pwd.Username,
		"user_id":  user.ID,
	})
	loginState, err = check.Check(ctx, loginState, CheckState{User: user, Password: pwd.Password})
	if err != nil {
		entry.WithError(err).Error("User password check failed")
		return loginState, err
	}

	entry.Debug("User password login successful")
	return loginState, nil
}

func (auth *facadeImpl) LoginUsingSecret(ctx context.Context, loginState LoginState, secret SecretLoginDTO) (LoginState, error) {
	// Find Entity

	// check secret

	// create tokens
	return nil, nil
}

func NewAuthFacade(usersRepo users.Repository, apps apps.Repository, jwkRepo jwtlib.JwkRepository) Facade {
	jwkService := jwtlib.NewJwkService(jwkRepo, usersRepo)
	jwtService := jwtlib.NewJwtService(jwkRepo, usersRepo, apps)
	return &facadeImpl{
		users:           usersRepo,
		passwordService: users.NewPasswordService(usersRepo),
		JwkService:      jwkService,
		JwtService:      jwtService,
	}
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
