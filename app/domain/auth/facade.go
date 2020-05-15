package auth

import (
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/domain/auth/jwtlib"
	"github.com/pestanko/gouthy/app/domain/entities"
	"github.com/pestanko/gouthy/app/domain/users"
	uuid "github.com/satori/go.uuid"
)

type Facade interface {
	LoginByID(id uuid.UUID) (Tokens, error)
}


type FacadeImpl struct {
	DB       *gorm.DB
	Users    users.Facade
	Entities entities.Facade
	Jwk      jwtlib.JwkInventory
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func (s *FacadeImpl) LoginByID(id uuid.UUID) (Tokens, error) {
	return Tokens{}, nil
}

func NewAuthFacade(db *gorm.DB, users users.Facade, entities entities.Facade, inventory jwtlib.JwkInventory) Facade {
	return &FacadeImpl{DB: db, Users: users, Entities: entities, Jwk: inventory}
}


