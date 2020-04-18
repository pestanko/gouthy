package services

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/models"
	"github.com/pestanko/gouthy/app/repositories"
	uuid "github.com/satori/go.uuid"
	"time"
)

type NewUser struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
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
	Username string    `json:"username"`
	Email    string    `json:"email"`
	ID       uuid.UUID `json:"id"`
}

type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
}

type UsersService interface {
	Create(newUser *NewUser) (models.User, error)
	Update(userId uuid.UUID, newUser *UpdateUser) (models.User, error)
	Delete(userId uuid.UUID) error
	UpdatePassword(userId uuid.UUID, password *UpdatePassword) error
	List() ([]ListUser, error)
	Get(userId uuid.UUID) (UserDTO, error)
	GetByUsername(userId string) (UserDTO, error)
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

func (s *UserServiceImpl) Get(id uuid.UUID) (UserDTO, error) {
	var user, err = s.users.FindByID(id)
	if err != nil {
		return UserDTO{}, err
	}

	return ConvertModelsToUserDTO(user), nil
}

func (s *UserServiceImpl) GetByUsername(username string) (UserDTO, error) {
	var user, err = s.users.FindByUsername(username)
	if err != nil {
		return UserDTO{}, err
	}

	return ConvertModelsToUserDTO(user), nil
}
