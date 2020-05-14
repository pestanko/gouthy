package auth

import (
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/domain/auth/jwtlib"
	"github.com/pestanko/gouthy/app/domain/entities"
	"github.com/pestanko/gouthy/app/domain/users"
	uuid "github.com/satori/go.uuid"
)

type Facade struct {
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

func (s *Facade) LoginByID(id uuid.UUID) Tokens {
	return Tokens{}
}

func NewAuthService(db *gorm.DB, users users.Facade, entities entities.Facade, inventory jwtlib.JwkInventory) *Facade {
	return &Facade{DB: db, Users: users, Entities: entities, Jwk: inventory}
}


