package users

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/shared/repositories"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	List(ctx context.Context) ([]User, error)
	FindByIdentifier(ctx context.Context, id string) (*User, error)
}

type repositoryDB struct {
	DB     *gorm.DB
	common repositories.CommonRepositoryDB
}

func NewUsersRepositoryDB(db *gorm.DB) Repository {
	return &repositoryDB{DB: db, common: repositories.NewCommonRepositoryDB(db, "User")}
}

func (r *repositoryDB) Create(ctx context.Context, user *User) error {
	return r.common.Create(ctx, user)
}

func (r *repositoryDB) Update(ctx context.Context, user *User) error {
	return r.common.Update(ctx, user)
}

func (r *repositoryDB) Delete(ctx context.Context, user *User) error {
	return r.common.Delete(ctx, user)
}

func (r *repositoryDB) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	result := r.DB.Where("id = ?", id).Find(&user)
	if result.Error != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"id": id,
		}).WithError(result.Error).Error("Find Failed")

		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, nil
		}

		return nil, result.Error
	}
	return &user, nil
}

func (r *repositoryDB) FindByIdentifier(ctx context.Context, id string) (*User, error) {
	uuidId, err := uuid.FromString(id)
	if err != nil {
		return r.FindByUsername(ctx, id)
	}
	return r.FindByID(ctx, uuidId)
}

func (r *repositoryDB) List(ctx context.Context) (result []User, err error) {
	r.DB.Find(&result)
	if r.DB.Error != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
		}).WithError(err).Error("List users failed")
	}

	return result, r.DB.Error
}

func (r *repositoryDB) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	result := r.DB.Where("username = ?", username).Find(&user)
	if result.Error != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"username": username,
		}).WithError(result.Error).Error("Find Failed")

		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, nil
		}

		return nil, result.Error
	}
	shared.GetLogger(ctx).WithFields(log.Fields{
		"username": username,
		"user_id":  user.ID,
	}).Debug("Found user")
	return &user, nil
}
