package jwtlib

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/pestanko/gouthy/app/shared"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"path"
)

type JwkRsa struct {
	algo       string
	keyID      string
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func (jwk *JwkRsa) PublicKey() crypto.PublicKey {
	return jwk.publicKey
}

func (jwk *JwkRsa) PrivateKey() crypto.PrivateKey {
	return jwk.privateKey
}

func (jwk *JwkRsa) Algorithm() string {
	return jwk.algo
}

func (jwk *JwkRsa) KeyId() string {
	return jwk.keyID
}

func NewJwkRsa(id string, algo string, pubKey *rsa.PublicKey, privKey *rsa.PrivateKey) Jwk {
	return &JwkRsa{keyID: id, algo: algo, publicKey: pubKey, privateKey: privKey}
}

func GenerateAndStoreNewRsaKey(ctx context.Context, basePath, name string) error {
	pthPrivateKey := path.Join(basePath, name+".pem")
	pthPublicKey := path.Join(basePath, name+".pub")
	bitSize := 4096

	shared.GetLogger(ctx).WithFields(log.Fields{
		"privatePath": pthPrivateKey,
		"publicPath":  pthPublicKey,
		"keySize":     bitSize,
	}).Debug("Generating RSA Keys")

	privateKey, err := generateRsaPrivateKey(ctx, bitSize)
	if err != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"privatePath": pthPrivateKey,
			"publicPath":  pthPublicKey,
		}).WithError(err).Error("Generating RSA Keys failed")
		return err
	}

	publicKeyBytes, err := encodeRsaPublicKey(ctx, &privateKey.PublicKey)
	if err != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"privatePath": pthPrivateKey,
			"publicPath":  pthPublicKey,
		}).WithError(err).Error("Encoding RSA public key failed")
		return err
	}

	privateKeyBytes := encodeRsaPrivateKeyToPEM(ctx, privateKey)

	err = writeKeyToFile(ctx, privateKeyBytes, pthPrivateKey)
	if err != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"privatePath": pthPrivateKey,
			"publicPath":  pthPublicKey,
		}).WithError(err).Error("Encoding RSA private key failed")
		return err
	}

	err = writeKeyToFile(ctx, publicKeyBytes, pthPublicKey)
	if err != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"privatePath": pthPrivateKey,
			"publicPath":  pthPublicKey,
		}).WithError(err).Error("Writing RSA key files failed")
		return err
	}
	return nil
}

func LoadRSAKey(ctx context.Context, basePath string, id string) (Jwk, error) {
	fp := path.Join(basePath, id)
	privatePath := fp + ".pem"
	publicPath := fp + ".pub"
	shared.GetLogger(ctx).WithFields(log.Fields{
		"privatePath": privatePath,
		"publicPath":  publicPath,
	}).Debug("Loading RSA Keys")

	contentPrivate, err := ioutil.ReadFile(privatePath)
	if err != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"privatePath": privatePath,
		}).WithError(err).Error("Unable to read a file")
		return nil, err
	}

	privateKey, err := parseRsaPrivateKeyFromPemStr(contentPrivate)
	if err != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"privatePath": privatePath,
		}).WithError(err).Error("Unable to parse a private key")
		return nil, err
	}

	contentPublic, err := ioutil.ReadFile(publicPath)
	if err != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"publicPath": publicPath,
		}).WithError(err).Error("Unable to read a file")
		return nil, err
	}

	pubKey, err := parseRsaPublicKeyFromPemStr(contentPublic)
	if err != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"publicPath": publicPath,
		}).WithError(err).Error("Unable to parse a public key")
		return nil, err
	}

	return NewJwkRsa(id, "RSA", pubKey, privateKey), nil
}

// generateRsaPrivateKey creates a RSA Private Key of specified byte size
func generateRsaPrivateKey(ctx context.Context, bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	shared.GetLogger(ctx).Info("Private Key generated")
	return privateKey, nil
}

// encodeRsaPrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodeRsaPrivateKeyToPEM(ctx context.Context, privateKey crypto.PrivateKey) []byte {
	// Get ASN.1 DER format
	keyBytes := x509.MarshalPKCS1PrivateKey(privateKey.(*rsa.PrivateKey))

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   keyBytes,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// encodeRsaPublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func encodeRsaPublicKey(ctx context.Context, publicKey crypto.PublicKey) ([]byte, error) {
	pk := publicKey.(*rsa.PublicKey)
	keyBytes := x509.MarshalPKCS1PublicKey(pk)

	var pemKeyBlock = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: keyBytes,
	}

	return pem.EncodeToMemory(pemKeyBlock), nil
}

// writePemToFile writes keys to a file
func writeKeyToFile(ctx context.Context, keyBytes []byte, filepath string) error {
	err := ioutil.WriteFile(filepath, keyBytes, 0600)
	if err != nil {
		return err
	}

	shared.GetLogger(ctx).WithField("filepath", filepath).Info("Key saved RSA KEY")
	return nil
}

func parseRsaPrivateKeyFromPemStr(privPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func parseRsaPublicKeyFromPemStr(pubPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pub, nil
}
