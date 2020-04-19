package services

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/models"
	"github.com/pestanko/gouthy/app/repositories"
	uuid "github.com/satori/go.uuid"
	"time"
)

type UserBase struct {
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	ID       uuid.UUID `json:"id"`
}

type NewUser struct {
	UserBase
	Password string `json:"password"`
}

type UpdateUser struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type UpdatePassword struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ListUser struct {
	UserBase
}

type UserDTO struct {
	UserBase
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UsersService interface {
	Create(newUser *NewUser) (models.User, error)
	Update(userId uuid.UUID, newUser *UpdateUser) (models.User, error)
	Delete(userId uuid.UUID) error
	UpdatePassword(userId uuid.UUID, password *UpdatePassword) error
	List() ([]ListUser, error)
	Get(userId uuid.UUID) (*models.User, error)
	GetByUsername(userId string) (*models.User, error)
	GetByAnyId(sid string) (*models.User, error)
}

type UserServiceImpl struct {
	users    repositories.UsersRepository
	entities repositories.EntitiesRepository
	secrets  repositories.SecretsService
}

func NewUsersService(db *gorm.DB) UsersService {
	return &UserServiceImpl{
		users:    repositories.NewUsersRepository(db),
		entities: repositories.NewEntitiesRepository(db),
		secrets:  repositories.NewSecretsRepository(db),
	}
}

func (s *UserServiceImpl) Create(newUser *NewUser) (models.User, error) {
	var entity = models.NewEntity()

	if err := s.entities.Create(entity); err != nil {
		return models.User{}, err
	}

	var user = models.User{
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

func (s *UserServiceImpl) Update(id uuid.UUID, update *UpdateUser) (models.User, error) {
	var user = models.User{
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

func (s *UserServiceImpl) UpdatePassword(id uuid.UUID, password *UpdatePassword) error {
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

func (s *UserServiceImpl) Delete(userId uuid.UUID) error {
	var user, err = s.users.FindByID(userId)
	if err != nil {
		return err
	}

	return s.users.Delete(user)
}

func (s *UserServiceImpl) List() ([]ListUser, error) {
	list, err := s.users.List()
	if err != nil {
		return []ListUser{}, err
	}

	listUsers := ConvertModelsToUserList(list)

	return listUsers, err
}

func (s *UserServiceImpl) Get(id uuid.UUID) (*models.User, error) {
	var user, err = s.users.FindByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) GetByUsername(username string) (*models.User, error) {
	var user, err = s.users.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) GetByAnyId(sid string) (*models.User, error) {
	var uid, err = uuid.FromString(sid)
	if err == nil {
		return s.Get(uid)
	}
	return s.GetByUsername(sid)
}
