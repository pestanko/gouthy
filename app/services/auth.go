package services

import (
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/jwtl"
	uuid "github.com/satori/go.uuid"
)

type AuthService struct {
	DB       *gorm.DB
	Users    UsersService
	Entities EntitiesService
	Jwk 	jwtl.JwkInventory
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func (s *AuthService) LoginByID(id uuid.UUID) AuthTokens {
	return AuthTokens{}
}

func NewAuthService(db *gorm.DB, users UsersService, entities EntitiesService, inventory jwtl.JwkInventory) *AuthService {
	return &AuthService{DB: db, Users: users, Entities: entities, Jwk: inventory}
}


