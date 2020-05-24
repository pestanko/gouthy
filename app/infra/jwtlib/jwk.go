package jwtlib

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/pestanko/gouthy/app/shared"
	log "github.com/sirupsen/logrus"
	jose "github.com/square/go-jose/v3"
	"io/ioutil"
	"path"
)

// Tutorial: https://medium.com/@niceoneallround/jwts-in-go-golang-4e0151f899af
// https://github.com/square/go-jose/blob/master/jose-util/utils.go

type JwkRepository interface {
	GetJwkForApplication(ctx context.Context, application string) (JwkSet, error)
	GetSystemJwk(ctx context.Context) (JwkSet, error)
	ListJwks(ctx context.Context) ([]string, error)
}

type JwkRepositoryImpl struct {
	KeysPath string
}

func (inventory JwkRepositoryImpl) ListJwks(ctx context.Context) (result []string, err error) {
	files, err := ioutil.ReadDir(inventory.KeysPath)
	if err != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"directory": inventory.KeysPath,
		}).WithError(err).Error("Unable to read directory")
		return result, err
	}

	for _, file := range files {
		if file.IsDir() {
			result = append(result, file.Name())
		}
	}

	return result, nil
}

func (inventory JwkRepositoryImpl) GetSystemJwk(ctx context.Context) (JwkSet, error) {
	return inventory.GetJwkForApplication(ctx, "system")
}

func (inventory JwkRepositoryImpl) GetJwkForApplication(ctx context.Context, appName string) (JwkSet, error) {
	return NewJwkSet(ctx, path.Join(inventory.KeysPath, appName))
}

func NewJwkRepository(keysPath string) JwkRepository {
	return &JwkRepositoryImpl{KeysPath: keysPath}
}

type JwkSet interface {
	GetById(keyID string) []Jwk
	GetOneById(keyId string) Jwk
}

type JwkSetImpl struct {
	set jose.JSONWebKeySet
}

func (set *JwkSetImpl) GetById(keyID string) []Jwk {
	return NewJwks(set.set.Key(keyID))
}

func (set *JwkSetImpl) GetOneById(keyID string) Jwk {
	keys := set.GetById(keyID)
	if len(keys) == 0 {
		return nil
	}
	return keys[0]
}

func NewJwks(keys []jose.JSONWebKey) (result []Jwk) {
	for _, k := range keys {
		result = append(result, NewJwk(k))
	}
	return result
}

func NewJwk(key jose.JSONWebKey) Jwk {
	return &JwkJose{key}
}

func NewJwkSet(ctx context.Context, pth string) (JwkSet, error) {
	files, err := ioutil.ReadDir(pth)
	if err != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"directory": pth,
		}).WithError(err).Error("Unable to read directory")
		return nil, err
	}
	var set jose.JSONWebKeySet
	for _, file := range files {
		fullPath := path.Join(pth, file.Name())
		content, err := ioutil.ReadFile(fullPath)
		if err != nil {
			shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
				"fullPath": fullPath,
				"name":     file.Name(),
			}).Warning("Unable to read JWK")
			continue
		}
		key, err := LoadPrivateKey(ctx, content)
		if err != nil {
			shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
				"fullPath": fullPath,
				"name":     file.Name(),
			}).Warning("Unable to parse JWK")
			continue
		}

		set.Keys = append(set.Keys, key.(jose.JSONWebKey))
	}
	return &JwkSetImpl{set: set}, nil
}

type Jwk interface {
	Algorithm() string
	KeyId() string
}

type JwkJose struct {
	key jose.JSONWebKey
}

func (jwk *JwkJose) Algorithm() string {
	return jwk.key.Algorithm
}

func (jwk *JwkJose) KeyId() string {
	return jwk.key.KeyID
}

// LoadPrivateKey loads a private key from PEM/DER/JWK-encoded data.
func LoadPrivateKey(ctx context.Context, data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	var priv interface{}
	priv, err := x509.ParsePKCS1PrivateKey(input)
	if err == nil {
		shared.GetLogger(ctx).WithError(err).Warning("Unable to parse x509 PKCS1 Private key")
		return priv, nil
	}

	priv, err = x509.ParsePKCS8PrivateKey(input)
	if err == nil {
		shared.GetLogger(ctx).WithError(err).Warning("Unable to parse x509 PKCS8 Private key")
		return priv, nil
	}

	priv, err = x509.ParseECPrivateKey(input)
	if err == nil {
		shared.GetLogger(ctx).WithError(err).Warning("Unable to parse x509 EC Private key")
		return priv, nil
	}

	jwk, err := loadJSONWebKey(ctx, input, false)
	if err == nil {
		shared.GetLogger(ctx).WithError(err).Warning("Unable to parse JSON Web Key Private key")
		return jwk, nil
	}

	return nil, errors.New("parse error, invalid private key")
}

func loadJSONWebKey(ctx context.Context, json []byte, pub bool) (*jose.JSONWebKey, error) {
	var jwk jose.JSONWebKey
	err := jwk.UnmarshalJSON(json)
	if err != nil {
		return nil, err
	}
	if !jwk.Valid() {
		return nil, errors.New("invalid JWK key")
	}
	if jwk.IsPublic() != pub {
		return nil, errors.New("priv/pub JWK key mismatch")
	}
	return &jwk, nil
}
