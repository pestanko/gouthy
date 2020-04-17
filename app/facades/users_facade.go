package facades

import (
	"fmt"
	"github.com/pestanko/gouthy/app/models"
	"github.com/pestanko/gouthy/app/services"
	uuid "github.com/satori/go.uuid"
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

type UsersFacade interface {
	Create(newUser *NewUser) (models.User, error)
	Update(userId uuid.UUID, newUser *UpdateUser) (models.User, error)
	Delete(userId uuid.UUID) error
	UpdatePassword(userId uuid.UUID, password UpdatePassword) error
	List() ([]ListUser, error)
	Get(userId uuid.UUID) error
}

type UserFacadeInteractor struct {
	userService    services.UsersService
	entityService  services.EntitiesService
	secretsService services.SecretsService
}

func (facade *UserFacadeInteractor) Create(newUser *NewUser) (models.User, error) {
	var entity = models.NewEntity()

	if err := facade.entityService.Create(entity); err != nil {
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

	if err := facade.userService.Create(&user); err != nil {
		return user, err
	}

	return user, nil
}

func (facade *UserFacadeInteractor) Update(userId uuid.UUID, update *UpdateUser) (models.User, error) {
	var user = models.User{
		Username: update.Username,
		Name:     update.Name,
		Email:    update.Email,
	}

	if err := facade.userService.Create(&user); err != nil {
		return user, err
	}

	return user, nil
}

func (facade *UserFacadeInteractor) UpdatePassword(userId uuid.UUID, password *UpdatePassword) error {
	var user, err = facade.userService.FindByID(userId)
	if err != nil {
		return err
	}

	if user.CheckPassword(password.CurrentPassword) {
		return fmt.Errorf("current password does not match")
	}

	if err = user.SetPassword(password.NewPassword); err != nil {
		return err
	}

	if err = facade.userService.Update(user); err != nil {
		return err
	}

	return nil
}

func (facade *UserFacadeInteractor) Delete(userId uuid.UUID) error {
	var user, err = facade.userService.FindByID(userId)
	if err != nil {
		return err
	}

	if err = facade.userService.Delete(user); err != nil {
		return err
	}

	return err
}
