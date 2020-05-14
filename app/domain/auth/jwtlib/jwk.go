package jwtlib

import (
	"path"
)

// Tutorial: https://medium.com/@niceoneallround/jwts-in-go-golang-4e0151f899af

type JwkInventory interface {
	GetJwkForRealm(realm string) (JwkSet, error)

	GetSystemJwk() (JwkSet, error)
}

type JwkInventoryImpl struct {
	KeysPath string
}

func (inventory JwkInventoryImpl) GetSystemJwk() (JwkSet, error) {
	return inventory.GetJwkForRealm("system")
}

func NewJwkInventory(keysPath string) JwkInventory {
	return &JwkInventoryImpl{KeysPath: keysPath}
}

func (inventory JwkInventoryImpl) GetJwkForRealm(realm string) (JwkSet, error) {
	return NewJwkSet(path.Join(inventory.KeysPath, realm)), nil
}

type JwkSet interface {
	GetById(keyID string) (Jwk, error)
}

type JwkSetImpl struct {
	path string
}

func (j JwkSetImpl) GetById(keyID string) (Jwk, error) {
	return nil, nil
}

func NewJwkSet(path string) JwkSet {
	return &JwkSetImpl{path: path}
}

type Jwk interface {
	GetPrivateKey()
	GetPublicKey()
}

type JwkImpl struct {

}
