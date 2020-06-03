package users

import (
	"context"
	"fmt"
	"github.com/pestanko/gouthy/app/shared"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type ListParams struct {
}

type Facade interface {
	Create(ctx context.Context, newUser *CreateDTO) (*UserDTO, error)
	Update(ctx context.Context, userId uuid.UUID, newUser *UpdateDTO) (*UserDTO, error)
	Delete(ctx context.Context, userId uuid.UUID) error
	UpdatePassword(ctx context.Context, userId uuid.UUID, password *UpdatePasswordDTO) error
	List(ctx context.Context, listParams ListParams) ([]ListUserDTO, error)
	Get(ctx context.Context, userId uuid.UUID) (*UserDTO, error)
	GetByUsername(ctx context.Context, userId string) (*UserDTO, error)
	GetByAnyId(ctx context.Context, sid string) (*UserDTO, error)
}

type facadeImpl struct {
	users   Repository
	secrets SecretsRepository
}

func NewUsersFacade(users Repository, secrets SecretsRepository) Facade {
	return &facadeImpl{users: users, secrets: secrets}
}

func (f *facadeImpl) Create(ctx context.Context, newUser *CreateDTO) (*UserDTO, error) {

	var user = &User{
		Username: newUser.Username,
		Name:     newUser.Name,
		Email:    newUser.Email,
	}

	if err := user.SetPassword(newUser.Password); err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"username": user.Username,
		}).Error("Unable to hash a password")
		return nil, err
	}

	if err := f.users.Create(ctx, user); err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"username": user.Username,
		}).Error("Unable to create a user")
		return nil, err
	}

	shared.GetLogger(ctx).WithFields(log.Fields{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("Creating a new user")

	return ConvertModelToDTO(user), nil
}

func (f *facadeImpl) Update(ctx context.Context, id uuid.UUID, update *UpdateDTO) (*UserDTO, error) {
	var user = User{
		Username: update.Username,
		Name:     update.Name,
		Email:    update.Email,
		ID:       id,
	}

	if err := f.users.Update(ctx, &user); err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"user_id":  user.ID,
			"username": user.Username,
		}).Error("Unable to update a user")
		return nil, err
	}

	return ConvertModelToDTO(&user), nil
}

func (f *facadeImpl) UpdatePassword(ctx context.Context, id uuid.UUID, password *UpdatePasswordDTO) error {
	var user, err = f.users.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if password.CurrentPassword != "" && user.CheckPassword(password.CurrentPassword) {
		return fmt.Errorf("current password does not match")
	}

	if err = user.SetPassword(password.NewPassword); err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"user_id":  user.ID,
			"username": user.Username,
		}).Error("Unable to hash a password")
		return err
	}

	return f.users.Update(ctx, user)
}

func (f *facadeImpl) Delete(ctx context.Context, userId uuid.UUID) error {
	var user, err = f.users.FindByID(ctx, userId)
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"user_id":  user.ID,
			"username": user.Username,
		}).Error("Unable to delete a user")
		return err
	}

	return f.users.Delete(ctx, user)
}

func (f *facadeImpl) List(ctx context.Context, listParams ListParams) ([]ListUserDTO, error) {
	list, err := f.users.List(ctx)
	if err != nil {
		return []ListUserDTO{}, err
	}

	return ConvertModelsToList(list), err
}

func (f *facadeImpl) Get(ctx context.Context, id uuid.UUID) (*UserDTO, error) {
	var user, err = f.users.FindByID(ctx, id)
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"user_id": id,
		}).Error("Unable to get a user")
		return nil, err
	}

	return ConvertModelToDTO(user), nil
}

func (f *facadeImpl) GetByUsername(ctx context.Context, username string) (*UserDTO, error) {
	var user, err = f.users.FindByUsername(ctx, username)
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"username": username,
		}).Error("Unable to get a user")
		return nil, err
	}

	return ConvertModelToDTO(user), nil
}

func (f *facadeImpl) GetByAnyId(ctx context.Context, sid string) (*UserDTO, error) {
	var uid, err = uuid.FromString(sid)
	if err == nil {
		return f.Get(ctx, uid)
	}

	return f.GetByUsername(ctx, sid)
}
