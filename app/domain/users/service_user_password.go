package users

import (
	"context"
	"github.com/pestanko/gouthy/app/shared"
	log "github.com/sirupsen/logrus"
)

type PasswordService interface {
	CheckPassword(ctx context.Context, user *User, provided string) (bool, error)
	SetPassword(ctx context.Context, user *User, provided string) error
}

type PasswordServiceImpl struct {
	repo Repository
}

func (s *PasswordServiceImpl) CheckPassword(ctx context.Context, user *User, provided string) (bool, error) {
	shared.GetLogger(ctx).WithFields(log.Fields{
		"user_id":       user.ID,
		"username": user.Username,
	}).Debug("Checking password")
	return user.CheckPassword(provided), nil
}

func (s *PasswordServiceImpl) SetPassword(ctx context.Context, user *User, provided string) error {
	if err := user.SetPassword(provided); err != nil {
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

func NewPasswordService(usersRepo Repository) PasswordService {
	return &PasswordServiceImpl{
		repo: usersRepo,
	}
}
