package auth

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/pestanko/gouthy/app/domain/apps"
	"github.com/pestanko/gouthy/app/domain/jwtlib"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/shared"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type Facade interface {
	Login(ctx context.Context, credentials Credentials) (LoginState, error)

	CreateSignedTokensFromLoginIdentity(ctx context.Context, identity *LoginIdentity) (SignedTokensDTO, error)
	CreateLoginIdentityFromToken(ctx context.Context, token jwtlib.Jwt) (*LoginIdentity, error)

	ParseJwt(ctx context.Context, str string) (jwtlib.Jwt, error)
	ParseAndValidateJwt(ctx context.Context, str string) (jwtlib.Jwt, error)
}

func NewAuthFacade(findUsers users.FindService, findApps apps.FindService, jwt jwtlib.JwtService, jwk jwtlib.JwkService, passwdService users.PasswordService) Facade {
	return &facadeImpl{
		passwordService: passwdService,
		JwkService:      jwk,
		JwtService:      jwt,
		findUsers:       findUsers,
		findApps:        findApps,
	}
}

type SignedTokensDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type Credentials struct {
	Username string
	Password string
	Totp     string
	Secret   string
}

func (d *SignedTokensDTO) Serialize() string {
	return shared.ToJSONIndent(d)
}

type facadeImpl struct {
	findUsers       users.FindService
	JwkService      jwtlib.JwkService
	JwtService      jwtlib.JwtService
	passwordService users.PasswordService
	findApps        apps.FindService
}

func (auth *facadeImpl) CreateSignedTokensFromLoginIdentity(ctx context.Context, identity *LoginIdentity) (SignedTokensDTO, error) {
	user, err := auth.findUsers.FindOne(ctx, users.FindQuery{AnyId: identity.UserId})
	if err != nil {
		return SignedTokensDTO{}, err
	}
	app, err := auth.findApps.FindOne(ctx, apps.FindQuery{AnyId: identity.ClientId})
	if err != nil {
		return SignedTokensDTO{}, err
	}

	return auth.CreateSignedTokensResponse(ctx, jwtlib.TokenCreateParams{
		User:   user,
		App:    app,
		Scopes: identity.Scopes,
	})
}

func (auth *facadeImpl) CreateLoginIdentityFromToken(ctx context.Context, token jwtlib.Jwt) (*LoginIdentity, error) {
	result := CreateLoginIdentityFromToken(token)

	return result, nil
}

func (auth *facadeImpl) ParseJwt(ctx context.Context, str string) (jwtlib.Jwt, error) {
	var claims jwt.MapClaims
	token, _, err := new(jwt.Parser).ParseUnverified(str, &claims)
	if err != nil {
		return nil, err
	}
	return jwtlib.NewJwt(token), nil
}

func (auth *facadeImpl) ParseAndValidateJwt(ctx context.Context, str string) (jwtlib.Jwt, error) {
	token, err := new(jwt.Parser).Parse(str, auth.getKeyId(ctx))
	if err != nil {
		return nil, err
	}
	return jwtlib.NewJwt(token), nil
}

func (auth *facadeImpl) getKeyId(ctx context.Context) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		id := token.Header["kid"]
		if id == "" {
			return nil, jwt.ErrInvalidKey
		}

		key, err := auth.JwkService.Get(ctx, id.(string))
		if err != nil {
			return nil, err
		}
		return key.PublicKey(), nil
	}
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

func (auth *facadeImpl) Login(ctx context.Context, cred Credentials) (LoginState, error) {
	user, loginState, err := auth.findUserForLoginState(ctx, cred.Username)

	if loginState != nil && loginState.IsNotOk() {
		return loginState, err
	}

	logEntry := shared.GetLogger(ctx).WithFields(log.Fields{
		"username": cred.Username,
		"user_id":  user.ID,
	})

	logEntry.Debug("Logging in the user")

	if cred.Password != "" {
		return auth.loginUsernamePassword(ctx, cred, loginState, user, logEntry)
	}

	return loginState, nil
}

func (auth *facadeImpl) findUserForLoginState(ctx context.Context, username string) (*users.User, LoginState, error) {
	user, err := auth.findUsers.FindOne(ctx, users.FindQuery{AnyId: username})

	if err != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"username": username,
		}).WithError(err).Error("Unable to find user - error happened")
		return nil, NewLoginState(uuid.UUID{}).AddStep(NewLoginStep(StepFindUser, Error)), err
	}
	if user == nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"username": username,
		}).Warning("Unable to find user")
		return nil, NewLoginState(uuid.UUID{}).AddStep(NewLoginStep(StepFindUser, Failed)), nil
	}

	shared.GetLogger(ctx).WithField("user", user).Debug("found user")

	return user, NewLoginState(user.ID).AddStep(NewLoginStep(StepFindUser, Success)), nil
}

func (auth *facadeImpl) loginUsernamePassword(ctx context.Context, cred Credentials, loginState LoginState, user *users.User, logEntry *log.Entry) (LoginState, error) {
	ch := NewLoginCheckPassword(auth.passwordService)
	loginState, err := ch.Check(ctx, loginState, CheckState{User: user, Password: cred.Password})
	if err != nil {
		logEntry.WithError(err).Error("User password check failed")
		return loginState.AddStep(NewLoginStep(StepLoginPassword, Error)), err
	}
	return loginState, nil
}

type TotpDTO struct {
	TotpCode string
}
