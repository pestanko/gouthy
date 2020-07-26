package jwtlib

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pestanko/gouthy/app/shared"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func NewJwt(token *jwt.Token) Jwt {
	impl := jwtImpl{
		token: token,
	}
	impl.jtiParts, _ = ParseJtiParts(impl.ID())
	return &impl
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
	Scopes() shared.Scopes
	Jti() JtiParts
	Raw() string

	// Kind of internal
	Claims() map[string]interface{}
	RawHeader() map[string]interface{}
}

type jwtImpl struct {
	token *jwt.Token
	jtiParts JtiParts
}

func (j *jwtImpl) RawHeader() map[string]interface{} {
	return j.token.Header
}

func (j *jwtImpl) Claims() map[string]interface{} {
	return j.token.Claims.(jwt.MapClaims)
}

func (j *jwtImpl) Raw() string {
	return j.token.Raw
}

func (j *jwtImpl) Jti() JtiParts {
	return j.jtiParts
}

func (j *jwtImpl) Scopes() shared.Scopes {
	scopes := j.mapClaims()["scope"].([]interface{})
	result := make(shared.Scopes, len(scopes))
	for i, item := range scopes {
		result[i] = item.(string)
	}
	return result
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
	shared.GetLogger(ctx).WithFields(log.Fields(serialize)).Info("Create a new TokenClaims")
	return NewJwt(token), nil
}

func NewJwtSigningService(repo JwkRepository) JwtSigningService {
	return &jwtSigningServiceImpl{keys: repo}
}


type JtiParts struct {
	CorrelationId string
	Type string
	Offset string
}

func ParseJtiParts(str string) (JtiParts, error) {
	if str == "" {
		return JtiParts{}, fmt.Errorf("jti string is empty")
	}
	parts := strings.Split(str, ".")
	if len(parts) != 3 {
		return JtiParts{}, fmt.Errorf("expected jti parts are 3, provided: %d", len(parts))

	}

	return JtiParts{
		Type:          parts[0],
		CorrelationId: parts[1],
		Offset:        parts[2],
	}, nil
}

func JtiPartsToString(parts JtiParts) string {
	return fmt.Sprintf("%s.%s.%s", parts.Type, parts.CorrelationId, parts.Offset)
}
