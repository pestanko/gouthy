package jwtlib

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	log "github.com/sirupsen/logrus"
	jose "github.com/square/go-jose/v3"
	"io/ioutil"
	"path"
)

// Tutorial: https://medium.com/@niceoneallround/jwts-in-go-golang-4e0151f899af
// https://github.com/square/go-jose/blob/master/jose-util/utils.go

type JwkRepository interface {
	GetJwkForRealm(realm string) (JwkSet, error)
	GetSystemJwk() (JwkSet, error)
	ListRealms() ([]string, error)
}

type JwkRepositoryImpl struct {
	KeysPath string
}

func (inventory JwkRepositoryImpl) ListRealms() (result []string, err error) {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		return result, err
	}

	for _, file := range files {
		if file.IsDir() {
			result = append(result, file.Name())
		}
	}

	return result, nil
}

func (inventory JwkRepositoryImpl) GetSystemJwk() (JwkSet, error) {
	return inventory.GetJwkForRealm("system")
}

func NewJwkInventory(keysPath string) JwkRepository {
	return &JwkRepositoryImpl{KeysPath: keysPath}
}

func (inventory JwkRepositoryImpl) GetJwkForRealm(realm string) (JwkSet, error) {
	return NewJwkSet(path.Join(inventory.KeysPath, realm))
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

func NewJwkSet(pth string) (JwkSet, error) {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		return nil, err
	}
	var set jose.JSONWebKeySet
	for _, file := range files {
		fullPath := path.Join(pth, file.Name())
		content, err := ioutil.ReadFile(fullPath)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"fullPath": fullPath,
				"name":     file.Name(),
			}).Error("Unable to read JWK")
			continue
		}
		key, err := LoadPrivateKey(content)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"fullPath": fullPath,
				"name":     file.Name(),
			}).Error("Unable to parse JWK")
			continue
		}

		set.Keys = append(set.Keys, key.(jose.JSONWebKey))
	}
	return &JwkSetImpl{set: set}, nil
}

type Jwk interface {
	GetPrivateKey()
	GetPublicKey()
}

type JwkJose struct {
	key jose.JSONWebKey
}

func (jwk *JwkJose) GetPrivateKey() {
	panic("implement me")
}

func (jwk *JwkJose) GetPublicKey() {
	panic("implement me")
}

// LoadPrivateKey loads a private key from PEM/DER/JWK-encoded data.
func LoadPrivateKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	var priv interface{}
	priv, err0 := x509.ParsePKCS1PrivateKey(input)
	if err0 == nil {
		return priv, nil
	}

	priv, err1 := x509.ParsePKCS8PrivateKey(input)
	if err1 == nil {
		return priv, nil
	}

	priv, err2 := x509.ParseECPrivateKey(input)
	if err2 == nil {
		return priv, nil
	}

	jwk, err3 := LoadJSONWebKey(input, false)
	if err3 == nil {
		return jwk, nil
	}

	return nil, errors.New("parse error, invalid private key")
}

func LoadJSONWebKey(json []byte, pub bool) (*jose.JSONWebKey, error) {
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

// LoadPublicKey loads a public key from PEM/DER/JWK-encoded data.
func LoadPublicKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	// Try to load SubjectPublicKeyInfo
	pub, err0 := x509.ParsePKIXPublicKey(input)
	if err0 == nil {
		return pub, nil
	}

	cert, err1 := x509.ParseCertificate(input)
	if err1 == nil {
		return cert.PublicKey, nil
	}

	jwk, err2 := LoadJSONWebKey(data, true)
	if err2 == nil {
		return jwk, nil
	}

	return nil, errors.New("parse error, invalid public key")
}
