package users

import (
	"context"
	"github.com/pestanko/gouthy/app/shared"
	log "github.com/sirupsen/logrus"
	"github.com/trustelem/zxcvbn"
)

func NewPasswordService(usersRepo Repository) PasswordService {
	return &PasswordServiceImpl{
		repo: usersRepo,
	}
}

type PasswordService interface {
	CheckPassword(ctx context.Context, user *User, provided string) (bool, error)
	SetPassword(ctx context.Context, user *User, provided string) error
}

type PasswordServiceImpl struct {
	repo Repository
}

func (s *PasswordServiceImpl) CheckPassword(ctx context.Context, user *User, provided string) (bool, error) {
	shared.GetLogger(ctx).WithFields(log.Fields{
		"user_id":  user.ID,
		"username": user.Username,
	}).Debug("Checking password")
	return user.CheckPassword(provided), nil
}

func (s *PasswordServiceImpl) SetPassword(ctx context.Context, user *User, password string) error {
	if err := user.SetPassword(password); err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"user_id":  user.ID,
			"username": user.Username,
		}).Error("Unable to hash a password")

		return err
	}

	if err := s.repo.Update(ctx, user); err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"user_id":  user.ID,
			"username": user.Username,
		}).Error("Unable to update a user")
		return err
	}
	return nil
}

type PasswordStrengthValidator struct {
	passwordPolicy *shared.FeatureConfig
}

func (v *PasswordStrengthValidator) FieldName() string {
	return "password"
}

func (v *PasswordStrengthValidator) Validate(data interface{}) shared.ValidationResult {
	result := shared.NewValidationResult()
	if !v.passwordPolicy.Enabled {
		return result
	}

	result = shared.NewFieldLengthValidator(v.getMin(), v.getMax(), "password").Validate(data)
	if result.IsFailed() {
		return result
	}

	strength := zxcvbn.PasswordStrength(data.(string), []string{})
	if strength.Score < v.getScore() {
		result.Fail("password score is less then required", shared.ValidationData{
			"expected":   v.getScore(),
			"provided":   strength.Score,
			"field_name": "password",
		})
	}

	return result
}

func (v *PasswordStrengthValidator) getMin() int {
	return v.passwordPolicy.GetParamInt("min_length", 0)
}

func (v *PasswordStrengthValidator) getMax() int {
	return v.passwordPolicy.GetParamInt("max_length", 0)
}

func (v *PasswordStrengthValidator) getScore() int {
	return v.passwordPolicy.GetParamInt("score", 0)
}

func NewPasswordValidator(passwordPolicy *shared.FeatureConfig) shared.Validator {
	return &PasswordStrengthValidator{
		passwordPolicy: passwordPolicy,
	}
}
