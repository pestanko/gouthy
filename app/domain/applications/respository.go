package applications

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/shared/repositories"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type Repository interface {
	Create(ctx context.Context, app *Application) error
	Update(ctx context.Context, app *Application) error
	Delete(ctx context.Context, app *Application) error
	FindByID(ctx context.Context, id uuid.UUID) (*Application, error)
	FindByIdentifier(ctx context.Context, id string) (*Application, error)
	FindByCodename(ctx context.Context, codename string) (*Application, error)
	List(ctx context.Context) ([]Application, error)
}

type repositoryDB struct {
	DB     *gorm.DB
	common repositories.CommonRepositoryDB
}

func (r *repositoryDB) FindByIdentifier(ctx context.Context, id string) (*Application, error) {
	uuidId, err := uuid.FromString(id)
	if err != nil {
		return r.FindByCodename(ctx, id)
	}
	return r.FindByID(ctx, uuidId)
}

func (r *repositoryDB) Create(ctx context.Context, user *Application) error {
	return r.common.Create(ctx, user)
}

func (r *repositoryDB) Update(ctx context.Context, user *Application) error {
	return r.common.Update(ctx, user)
}

func (r *repositoryDB) Delete(ctx context.Context, user *Application) error {
	return r.common.Delete(ctx, user)
}

func (r *repositoryDB) FindByID(ctx context.Context, id uuid.UUID) (*Application, error) {
	var app Application
	result := r.DB.Where("id = ?", id).Find(&app)
	if result.Error != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"id": id,
		}).WithError(result.Error).Error("Find Failed")

		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, nil
		}

		return nil, result.Error
	}
	return &app, nil
}

func (r *repositoryDB) FindByCodename(ctx context.Context, codename string) (*Application, error) {
	var application Application
	result := r.DB.Where("codename = ?", codename).Find(&application)
	if result.Error != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"codename": codename,
		}).WithError(result.Error).Error("Find Failed")

		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, nil
		}

		return nil, result.Error
	}
	shared.GetLogger(ctx).WithFields(log.Fields{
		"codename":       codename,
		"application_id": application.ID,
	}).Debug("Found application")
	return &application, nil
}

func (r *repositoryDB) List(ctx context.Context) (result []Application, err error) {
	r.DB.Find(&result)
	if r.DB.Error != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
		}).WithError(err).Error("List applications failed")
	}
	return result, r.DB.Error
}

func NweApplicationsRepositoryDB(db *gorm.DB) Repository {
	return &repositoryDB{
		DB:     db,
		common: repositories.NewCommonRepositoryDB(db, "Applications"),
	}
}
