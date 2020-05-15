package auth

import (
	"github.com/pestanko/gouthy/app/domain/auth/jwtlib"
	"github.com/pestanko/gouthy/app/domain/entities"
	"github.com/pestanko/gouthy/app/domain/machines"
	"github.com/pestanko/gouthy/app/domain/users"
	uuid "github.com/satori/go.uuid"
)

type Facade interface {
	LoginUsernamePassword(pwd PasswordLoginDTO) (Tokens, error)
	LoginUsingSecret(secret SecretLoginDTO) (Tokens, error)
}

type FacadeImpl struct {
	Users    users.Repository
	Machines machines.Repository
	Entities entities.Repository
	Jwk      jwtlib.JwkInventory
}

func (auth *FacadeImpl) LoginUsernamePassword(pwd PasswordLoginDTO) (Tokens, error) {
	user, err := auth.Users.FindByUsername(pwd.Username)
	if err != nil {
		return Tokens{}, err
	}
	flow := UcPasswordFlow{Users: auth.Users}

	if flow.Check(user, pwd.Password) != nil {
		return Tokens{}, err
	}
	return auth.createTokensForEntity(user.ID)
}

func (auth *FacadeImpl) LoginUsingSecret(secret SecretLoginDTO) (Tokens, error) {
	// Get Entity

	// check secret

	// create tokens
	return Tokens{}, nil
}

func (auth *FacadeImpl) createTokensForEntity(id uuid.UUID) (Tokens, error) {
	return Tokens{}, nil
}

func NewAuthFacade(users users.Repository, machines machines.Repository, entities entities.Repository, inventory jwtlib.JwkInventory) Facade {
	return &FacadeImpl{Users: users, Entities: entities, Jwk: inventory, Machines: machines}
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type PasswordLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SecretLoginDTO struct {
	Secret     string `json:"secret"`
	Codename   string `json:"codename"`
	EntityType string `json:"entity_type"`
}
