package jwtlib

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/pestanko/gouthy/app/shared"
	log "github.com/sirupsen/logrus"
	"time"
)

func NewJwt(token *jwt.Token) Jwt {
	return &jwtImpl{
		token: token,
	}
}

type Jwt interface {
	JwkID() string
	ClientId() string
	Audience() string
	Subject() string
	ID() string
	Issuer() string

	IssuedAt() time.Time
	ExpiresAt() time.Time
	UserId() string
	AppId() string
	Scopes() []string
}

type jwtImpl struct {
	token *jwt.Token
}

func (j *jwtImpl) Scopes() []string {
	return j.Scopes()
}

func (j *jwtImpl) AppId() string {
	return j.Audience()
}

func (j *jwtImpl) UserId() string {
	return j.Subject()
}

func (j *jwtImpl) Issuer() string {
	return j.stringClaim("iss")
}

func (j *jwtImpl) IssuedAt() time.Time {
	return j.timeClaim("iat")
}

func (j *jwtImpl) ExpiresAt() time.Time {
	return j.timeClaim("exp")
}

func (j *jwtImpl) Subject() string {
	return j.stringClaim("sub")
}

func (j *jwtImpl) ID() string {
	return j.stringClaim("jti")
}

func (j *jwtImpl) ClientId() string {
	return j.Audience()
}

func (j *jwtImpl) JwkID() string {
	return j.token.Header["kid"].(string)
}

func (j *jwtImpl) Audience() string {
	return j.stringClaim("aud")
}

func (j *jwtImpl) mapClaims() jwt.MapClaims {
	return j.token.Claims.(jwt.MapClaims)
}

func (j *jwtImpl) stringClaim(claim string) string {
	return j.mapClaims()[claim].(string)
}

func (j *jwtImpl) timeClaim(claim string) time.Time {
	switch iat := j.mapClaims()[claim].(type) {
	case float64:
		return time.Unix(int64(iat), 0)
	case json.Number:
		v, _ := iat.Int64()
		return time.Unix(v, 0)
	}
	return time.Time{}
}

type JwtSigningService interface {
	Sign(ctx context.Context, token Jwt) (*SignedJwt, error)
	Create(ctx context.Context, claims Claims) (Jwt, error)
}

type jwtSigningServiceImpl struct {
	keys JwkRepository
}

type SignedJwt struct {
	jwtImpl
	Signature string
}

func (service *jwtSigningServiceImpl) Sign(ctx context.Context, token Jwt) (*SignedJwt, error) {
	impl := token.(*jwtImpl)

	key, err := service.keys.GetLatest(ctx)

	if err != nil {
		return nil, err
	}

	impl.token.Header["kid"] = key.KeyId()

	signed, err := impl.token.SignedString(key.PrivateKey())
	if err != nil {
		return nil, err
	}
	return &SignedJwt{*impl, signed}, nil
}

func (service *jwtSigningServiceImpl) Create(ctx context.Context, claims Claims) (Jwt, error) {
	serialize := claims.Serialize()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), serialize)
	shared.GetLogger(ctx).WithFields(log.Fields(serialize)).Info("Create a new Jwt")
	return NewJwt(token), nil
}

func NewJwtSigningService(repo JwkRepository) JwtSigningService {
	return &jwtSigningServiceImpl{keys: repo}
}
