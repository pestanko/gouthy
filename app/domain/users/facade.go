package users

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/domain/entities"
	uuid "github.com/satori/go.uuid"
)



type Facade interface {
	Create(newUser *NewUserDTO) (User, error)
	Update(userId uuid.UUID, newUser *UpdateUserDTO) (User, error)
	Delete(userId uuid.UUID) error
	UpdatePassword(userId uuid.UUID, password *UpdatePasswordDTO) error
	List() ([]ListUserDTO, error)
	Get(userId uuid.UUID) (*User, error)
	GetByUsername(userId string) (*User, error)
	GetByAnyId(sid string) (*User, error)
}

type FacadeImpl struct {
	users    Repository
	entities entities.Repository
	secrets  entities.SecretsRepository
}

func NewUsersFacade(db *gorm.DB) Facade {
	return &FacadeImpl{
		users:    NewUsersRepository(db),
		entities: entities.NewEntitiesRepository(db),
		secrets:  entities.NewSecretsRepository(db),
	}
}

func (s *FacadeImpl) Create(newUser *NewUserDTO) (User, error) {
	var entity = entities.NewEntity()

	if err := s.entities.Create(entity); err != nil {
		return User{}, err
	}

	var user = User{
		ID:       entity.ID,
		Username: newUser.Username,
		Name:     newUser.Name,
		Email:    newUser.Email,
	}

	if err := user.SetPassword(newUser.Password); err != nil {
		return user, err
	}

	if err := s.users.Create(&user); err != nil {
		return user, err
	}

	return user, nil
}

func (s *FacadeImpl) Update(id uuid.UUID, update *UpdateUserDTO) (User, error) {
	var user = User{
		Username: update.Username,
		Name:     update.Name,
		Email:    update.Email,
		ID:       id,
	}

	if err := s.users.Update(&user); err != nil {
		return user, err
	}

	return user, nil
}

func (s *FacadeImpl) UpdatePassword(id uuid.UUID, password *UpdatePasswordDTO) error {
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

func (s *FacadeImpl) Delete(userId uuid.UUID) error {
	var user, err = s.users.FindByID(userId)
	if err != nil {
		return err
	}

	return s.users.Delete(user)
}

func (s *FacadeImpl) List() ([]ListUserDTO, error) {
	list, err := s.users.List()
	if err != nil {
		return []ListUserDTO{}, err
	}

	listUsers := ConvertModelsToUserList(list)

	return listUsers, err
}

func (s *FacadeImpl) Get(id uuid.UUID) (*User, error) {
	var user, err = s.users.FindByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *FacadeImpl) GetByUsername(username string) (*User, error) {
	var user, err = s.users.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *FacadeImpl) GetByAnyId(sid string) (*User, error) {
	var uid, err = uuid.FromString(sid)
	if err == nil {
		return s.Get(uid)
	}
	return s.GetByUsername(sid)
}
