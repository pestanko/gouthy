package users

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/domain/entities"
	uuid "github.com/satori/go.uuid"
)

type PasswordLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Facade interface {
	Create(newUser *NewUserDTO) (*UserDTO, error)
	Update(userId uuid.UUID, newUser *UpdateUserDTO) (*UserDTO, error)
	Delete(userId uuid.UUID) error
	UpdatePassword(userId uuid.UUID, password *UpdatePasswordDTO) error
	CheckPassword(pwdLogin PasswordLoginDTO) error
	List() ([]ListUserDTO, error)
	Get(userId uuid.UUID) (*UserDTO, error)
	GetByUsername(userId string) (*UserDTO, error)
	GetByAnyId(sid string) (*UserDTO, error)
}

type facadeImpl struct {
	users    Repository
	entities entities.Repository
	secrets  entities.SecretsRepository
}

func NewUsersFacade(db *gorm.DB) Facade {
	return &facadeImpl{
		users:    NewUsersRepository(db),
		entities: entities.NewEntitiesRepository(db),
		secrets:  entities.NewSecretsRepository(db),
	}
}

func (s *facadeImpl) Create(newUser *NewUserDTO) (*UserDTO, error) {
	var entity = entities.NewEntity()

	if err := s.entities.Create(entity); err != nil {
		return nil, err
	}

	var user = User{
		ID:       entity.ID,
		Username: newUser.Username,
		Name:     newUser.Name,
		Email:    newUser.Email,
	}

	if err := user.SetPassword(newUser.Password); err != nil {
		return nil, err
	}

	if err := s.users.Create(&user); err != nil {
		return nil, err
	}

	return ConvertModelToUserDTO(&user), nil
}

func (s *facadeImpl) Update(id uuid.UUID, update *UpdateUserDTO) (*UserDTO, error) {
	var user = User{
		Username: update.Username,
		Name:     update.Name,
		Email:    update.Email,
		ID:       id,
	}

	if err := s.users.Update(&user); err != nil {
		return nil, err
	}

	return ConvertModelToUserDTO(&user), nil
}

func (s *facadeImpl) CheckPassword(pwdLogin PasswordLoginDTO) error {
	user, err := s.users.FindByUsername(pwdLogin.Username)
	if err != nil {
		return err
	}

	if user.CheckPassword(pwdLogin.Password) {
		return fmt.Errorf("password does not match")
	}

	return nil
}

func (s *facadeImpl) UpdatePassword(id uuid.UUID, password *UpdatePasswordDTO) error {
	var user, err = s.users.FindByID(id)
	if err != nil {
		return err
	}

	if password.CurrentPassword != "" && user.CheckPassword(password.CurrentPassword) {
		return fmt.Errorf("current password does not match")
	}

	if err = user.SetPassword(password.NewPassword); err != nil {
		return err
	}

	return s.users.Update(user)
}

func (s *facadeImpl) Delete(userId uuid.UUID) error {
	var user, err = s.users.FindByID(userId)
	if err != nil {
		return err
	}

	return s.users.Delete(user)
}

func (s *facadeImpl) List() ([]ListUserDTO, error) {
	list, err := s.users.List()
	if err != nil {
		return []ListUserDTO{}, err
	}

	listUsers := ConvertModelsToUserList(list)

	return listUsers, err
}

func (s *facadeImpl) Get(id uuid.UUID) (*UserDTO, error) {
	var user, err = s.users.FindByID(id)
	if err != nil {
		return nil, err
	}

	return ConvertModelToUserDTO(user), nil
}

func (s *facadeImpl) GetByUsername(username string) (*UserDTO, error) {
	var user, err = s.users.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	return ConvertModelToUserDTO(user), nil
}

func (s *facadeImpl) GetByAnyId(sid string) (*UserDTO, error) {
	var uid, err = uuid.FromString(sid)
	if err == nil {
		return s.Get(uid)
	}

	return s.GetByUsername(sid)
}
