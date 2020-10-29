package users

import (
	"context"
	"fmt"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/shared/repos"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type ListParams struct {
	Limit  int
	Offset int
}

type Facade interface {
	Create(ctx context.Context, newUser *CreateDTO) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, userId uuid.UUID) error
	UpdatePassword(ctx context.Context, userId uuid.UUID, password *UpdatePasswordDTO) error
	List(ctx context.Context, params ListParams) ([]User, error)
	Get(ctx context.Context, userId uuid.UUID) (*User, error)
	GetByUsername(ctx context.Context, userId string) (*User, error)
	GetByAnyId(ctx context.Context, sid string) (*User, error)
}

type facadeImpl struct {
	users           Repository
	secrets         SecretsRepository
	findService     FindService
	passwordService PasswordService
	features        *shared.FeaturesConfig
}

func NewUsersFacade(users Repository, secrets SecretsRepository, services Services, features *shared.FeaturesConfig) Facade {
	return &facadeImpl{
		users:           users,
		secrets:         secrets,
		findService:     services.Find,
		passwordService: services.Password,
		features:        features,
	}
}

func (f *facadeImpl) Create(ctx context.Context, newUser *CreateDTO) (*User, error) {

	if err := f.checkCreateUserRequirements(newUser); err != nil {
		err.LogAppend(shared.GetLogger(ctx)).WithFields(newUser.LogFields()).Error("User validation failed")
		return nil, err
	}

	var user = &User{
		Username: newUser.Username,
		Name:     newUser.Name,
		Email:    newUser.Email,
	}

	if err := user.SetPassword(newUser.Password); err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(user.LogFields()).Error("Unable to hash a password")
		return nil, err
	}

	if err := f.users.Create(ctx, user); err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(user.LogFields()).Error("Unable to create a new user")
		return nil, err
	}

	shared.GetLogger(ctx).WithFields(user.LogFields()).Info("Creating a new user")

	return user, nil
}

func (f *facadeImpl) Update(ctx context.Context, user *User) (*User, error) {
	if err := f.users.Update(ctx, user); err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(user.LogFields()).Error("Unable to update a user")
		return nil, err
	}

	return user, nil
}

func (f *facadeImpl) UpdatePassword(ctx context.Context, id uuid.UUID, password *UpdatePasswordDTO) error {
	var user, err = f.users.QueryOne(ctx, FindQuery{Id: id})
	if err != nil {
		return err
	}

	if password.CurrentPassword != "" && user.CheckPassword(password.CurrentPassword) {
		shared.GetLogger(ctx).WithError(err).WithFields(user.LogFields()).Error("Unable change the password, current password does not match")
		return fmt.Errorf("current password does not match")
	}

	if err := f.checkPassword(password.NewPassword); err != nil {
		err.LogAppend(shared.GetLogger(ctx)).WithFields(user.LogFields()).Error("Unable change the password, password check failed")
		return err
	}

	return f.passwordService.SetPassword(ctx, user, password.NewPassword)
}

func (f *facadeImpl) Delete(ctx context.Context, userId uuid.UUID) error {
	var user, err = f.users.QueryOne(ctx, FindQuery{Id: userId})
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(user.LogFields()).Error("Unable to delete a user")
		return err
	}

	return f.users.Delete(ctx, user)
}

func (f *facadeImpl) List(ctx context.Context, params ListParams) ([]User, error) {
	list, err := f.findService.Find(ctx, FindQuery{
		PaginationQuery: repos.NewPaginationQuery(params.Limit, params.Offset),
	})
	if err != nil {
		return []User{}, err
	}

	return list, err
}

func (f *facadeImpl) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	var user, err = f.findService.FindOne(ctx, FindQuery{Id: id})
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"user_id": id,
		}).Error("Unable to get a user")
		return nil, err
	}

	return user, nil
}

func (f *facadeImpl) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user, err = f.findService.FindOne(ctx, FindQuery{Username: username})
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"username": username,
		}).Error("Unable to get a user")
		return nil, err
	}

	return user, nil
}

func (f *facadeImpl) GetByAnyId(ctx context.Context, sid string) (*User, error) {
	one, err := f.findService.FindOne(ctx, FindQuery{AnyId: sid})
	return one, err
}

func (f *facadeImpl) checkCreateUserRequirements(user *CreateDTO) shared.AppError {
	result := shared.NewFieldLengthValidator(3, 50, "username").Validate(user.Username)
	if result.IsFailed() {
		return result.IntoError()
	}

	result = shared.NewEmailValidator().Validate(user.Email)
	if result.IsFailed() {
		return result.IntoError()
	}

	result = NewPasswordValidator(f.features.PasswordPolicy()).Validate(user.Password)
	if result.IsFailed() {
		return result.IntoError()
	}

	return nil
}

func (f *facadeImpl) checkPassword(password string) shared.AppError {
	result := NewPasswordValidator(f.features.PasswordPolicy()).Validate(password)
	return shared.NewErrFromValidator(&result)
}
