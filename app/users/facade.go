package users

import (
	"context"
	"fmt"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/shared/repos"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/trustelem/zxcvbn"
	"regexp"
)

type ListParams struct {
	Limit  int
	Offset int
}

type Facade interface {
	Create(ctx context.Context, newUser *CreateDTO) (*UserDTO, error)
	Update(ctx context.Context, userId uuid.UUID, newUser *UpdateDTO) (*UserDTO, error)
	Delete(ctx context.Context, userId uuid.UUID) error
	UpdatePassword(ctx context.Context, userId uuid.UUID, password *UpdatePasswordDTO) error
	List(ctx context.Context, params ListParams) ([]ListUserDTO, error)
	Get(ctx context.Context, userId uuid.UUID) (*UserDTO, error)
	GetByUsername(ctx context.Context, userId string) (*UserDTO, error)
	GetByAnyId(ctx context.Context, sid string) (*UserDTO, error)
}

type facadeImpl struct {
	users           Repository
	secrets         SecretsRepository
	findService     FindService
	passwordService PasswordService
}

func NewUsersFacade(users Repository, secrets SecretsRepository, services Services) Facade {
	return &facadeImpl{
		users:           users,
		secrets:         secrets,
		findService:     services.Find,
		passwordService: services.Password,
	}
}

func (f *facadeImpl) Create(ctx context.Context, newUser *CreateDTO) (*UserDTO, error) {

	if err := checkCreateUserRequirements(newUser); err != nil {
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
		shared.GetLogger(ctx).WithError(err).WithFields(user.LogFields()).Error("Unable to update a user")
		return nil, err
	}

	return ConvertModelToDTO(&user), nil
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

	if err := checkPassword(password.NewPassword); err != nil {
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

func (f *facadeImpl) List(ctx context.Context, params ListParams) ([]ListUserDTO, error) {
	list, err := f.findService.Find(ctx, FindQuery{
		PaginationQuery: repos.NewPaginationQuery(params.Limit, params.Offset),
	})
	if err != nil {
		return []ListUserDTO{}, err
	}

	return ConvertModelsToList(list), err
}

func (f *facadeImpl) Get(ctx context.Context, id uuid.UUID) (*UserDTO, error) {
	var user, err = f.findService.FindOne(ctx, FindQuery{Id: id})
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"user_id": id,
		}).Error("Unable to get a user")
		return nil, err
	}

	return ConvertModelToDTO(user), nil
}

func (f *facadeImpl) GetByUsername(ctx context.Context, username string) (*UserDTO, error) {
	var user, err = f.findService.FindOne(ctx, FindQuery{Username: username})
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"username": username,
		}).Error("Unable to get a user")
		return nil, err
	}

	return ConvertModelToDTO(user), nil
}

func (f *facadeImpl) GetByAnyId(ctx context.Context, sid string) (*UserDTO, error) {
	one, err := f.findService.FindOne(ctx, FindQuery{AnyId: sid})
	return ConvertModelToDTO(one), err
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// checkValidEmail checks if the email provided passes the required structure and length.
func checkValidEmail(e string) shared.AppError {
	if err := shared.CheckFieldMinLength(e, "email", 3); err != nil {
		return err
	}

	if err := shared.CheckFieldMaxLength(e, "email", 254); err != nil {
		return err
	}

	if !emailRegex.MatchString(e) {
		return shared.NewErrInvalidField("email").WithDetail(shared.ErrDetail{
			"reason": "Provided email is not valid",
			"email":  e,
		})
	}
	return nil
}

func checkCreateUserRequirements(user *CreateDTO) shared.AppError {
	err := checkUsername(user.Username)
	if err != nil {
		return err
	}

	if err := checkValidEmail(user.Email); err != nil {
		return err
	}

	if err := checkPassword(user.Password); err != nil {
		return err
	}
	return nil
}

func checkUsername(username string) shared.AppError {
	if err := shared.CheckFieldMinLength(username, "username", 3); err != nil {
		return err
	}

	if err := shared.CheckFieldMaxLength(username, "username", 50); err != nil {
		return err
	}
	return nil
}

func checkPassword(password string) shared.AppError {
	if err := shared.CheckFieldMinLength(password, "password", 8); err != nil {
		return err
	}

	if err := shared.CheckFieldMaxLength(password, "password", 256); err != nil {
		return err
	}

	compl := zxcvbn.PasswordStrength(password, []string{})
	if compl.Score < 3 {
		return NewErrInvalidPassword().WithDetail(shared.ErrDetail{
			"reason":    "Password too weak",
			"score":     compl.Score,
			"min_score": 3,
			"calc_time": compl.CalcTime,
			"guesses":   compl.Guesses,
		})
	}

	return nil
}