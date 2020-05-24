package repositories

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/shared"
	log "github.com/sirupsen/logrus"
)

type CommonRepositoryDB struct {
	db   *gorm.DB
	Name string
}

func NewCommonRepositoryDB(db *gorm.DB, name string) CommonRepositoryDB {
	return CommonRepositoryDB{db: db, Name: name}
}

func (repo *CommonRepositoryDB) Create(ctx context.Context, instance interface{}) error {
	result := repo.db.Create(instance)
	if result.Error != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"repo":     repo.Name,
			"action":   "create",
			"instance": instance,
		}).WithError(result.Error).Error("Create failed")
		return result.Error
	}
	shared.GetLogger(ctx).WithFields(log.Fields{
		"repo":     repo.Name,
		"action":   "create",
		"instance": instance,
	}).Info("Created instance")
	return nil
}

func (repo *CommonRepositoryDB) Update(ctx context.Context, instance interface{}) error {
	result := repo.db.Update(instance)
	if result.Error != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"repo":     repo.Name,
			"action":   "update",
			"instance": instance,
		}).WithError(result.Error).Error("Update failed")
		return result.Error
	}
	shared.GetLogger(ctx).WithFields(log.Fields{
		"repo":     repo.Name,
		"action":   "update",
		"instance": instance,
	}).Info("Updated instance")
	return nil
}

func (repo *CommonRepositoryDB) Delete(ctx context.Context, instance interface{}) error {
	result := repo.db.Delete(instance)
	if result.Error != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"repo":     repo.Name,
			"action":   "delete",
			"instance": instance,
		}).WithError(result.Error).Error("Delete failed")
		return result.Error
	}
	shared.GetLogger(ctx).WithFields(log.Fields{
		"repo":     repo.Name,
		"action":   "delete",
		"instance": instance,
	}).Info("Deleted instance")
	return nil
}
