package jwtlib

import (
	"context"
	"crypto"
	"fmt"
	"github.com/google/uuid"
	"github.com/pestanko/gouthy/app/shared"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Tutorial: https://medium.com/@niceoneallround/jwts-in-go-golang-4e0151f899af
// https://github.com/square/go-jose/blob/master/jose-util/utils.go

type JwkGenerateParams struct {
}

type JwkRepository interface {
	Get(ctx context.Context, id string) (Jwk, error)
	List(ctx context.Context) ([]Jwk, error)
	Generate(ctx context.Context, params JwkGenerateParams) error
	Store(ctx context.Context, jwk Jwk) error
	Add(ctx context.Context, jwk Jwk) error
}

type JwkRepositoryImpl struct {
	BasePath string
	cache    map[string]Jwk
}

func (repo *JwkRepositoryImpl) Store(ctx context.Context, jwk Jwk) error {
	pthPrivateKey := path.Join(repo.BasePath, jwk.KeyId()+".pem")
	pthPublicKey := path.Join(repo.BasePath, jwk.KeyId()+".pub")

	publicKeyBytes, err := encodeRsaPublicKey(ctx, jwk.PublicKey())
	if err != nil {
		return err
	}

	privateKeyBytes := encodeRsaPrivateKeyToPEM(ctx, jwk.PrivateKey())

	err = writeKeyToFile(ctx, privateKeyBytes, pthPrivateKey)
	if err != nil {
		return err
	}

	err = writeKeyToFile(ctx, publicKeyBytes, pthPublicKey)
	if err != nil {
		return err
	}
	return nil
}

func (repo *JwkRepositoryImpl) Generate(ctx context.Context, params JwkGenerateParams) error {
	// TODO: Better implementation
	keyId, err := uuid.NewRandom()
	datetime := time.Now()
	strTime := datetime.Format("2006-01-02T15-04-05")
	name := fmt.Sprintf("%s_%s", strTime, keyId.String())
	if err != nil {
		return err
	}
	return GenerateAndStoreNewRsaKey(ctx, repo.BasePath, name)
}

func (repo *JwkRepositoryImpl) Add(ctx context.Context, jwk Jwk) error {
	shared.GetLogger(ctx).WithFields(log.Fields{
		"id":  jwk.KeyId(),
		"alg": jwk.Algorithm(),
	}).Info("Adding the jwk to the repository")
	repo.cache[jwk.KeyId()] = jwk
	return nil
}

func (repo *JwkRepositoryImpl) Get(ctx context.Context, id string) (Jwk, error) {
	val, ok := repo.cache[id]
	if ok {
		return val, nil
	}

	fpath := path.Join(repo.BasePath, id)
	if !Exists(fpath) {
		return nil, fmt.Errorf("key not found - %s", id)

	}

	key, err := LoadRSAKey(ctx, repo.BasePath, id)
	if err != nil {
		return nil, err
	}

	if err := repo.Add(ctx, key); err != nil {
		return nil, err
	}

	return key, nil
}

func (repo *JwkRepositoryImpl) List(ctx context.Context) (result []Jwk, err error) {
	if len(repo.cache) == 0 {
		if err = repo.loadCache(ctx); err != nil {
			return nil, err
		}
	}
	for _, key := range repo.cache {
		result = append(result, key)
	}
	return result, nil
}

func (repo *JwkRepositoryImpl) loadCache(ctx context.Context) error {
	shared.GetLogger(ctx).WithField("path", repo.BasePath).Info("Loading cache")
	files, err := filepath.Glob(path.Join(repo.BasePath, "*.pem"))
	if err != nil {
		return err
	}
	for _, filePath := range files {
		fullname := path.Base(filePath)
		name := strings.TrimSuffix(fullname, ".pem")
		log.WithFields(log.Fields{
			"filePath": filePath,
			"name":     name,
		}).Debug("Found key")

		key, err := LoadRSAKey(ctx, repo.BasePath, name)
		if err != nil {
			return err
		}

		if err := repo.Add(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

func NewJwkRepository(keysPath string) JwkRepository {
	return &JwkRepositoryImpl{BasePath: keysPath, cache: make(map[string]Jwk)}
}

type Jwk interface {
	Algorithm() string
	KeyId() string

	PublicKey() crypto.PublicKey
	PrivateKey() crypto.PrivateKey
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
