package repositories

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/shared"
	log "github.com/sirupsen/logrus"
)

type PaginationQuery struct {
	Offset int
	Limit  int
}

func NewPaginationQuery(limit, offset int) PaginationQuery {
	return PaginationQuery{Limit: limit, Offset: offset}
}

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

func (repo *CommonRepositoryDB) ProcessQueryOne(db *gorm.DB, result interface{}, entry *log.Entry) (interface{}, error) {
	dbResult := db.First(result)
	if dbResult.Error != nil {
		if gorm.IsRecordNotFoundError(dbResult.Error) {
			entry.Debug("Record not found")
			return nil, nil
		}
		entry.WithError(dbResult.Error).Error("Unable to query the record")
		return nil, dbResult.Error
	}
	entry.Trace("Query the record")
	return result, nil
}

func (repo *CommonRepositoryDB) ProcessQuery(db *gorm.DB, result interface{}, entry *log.Entry) error {
	db.Find(result)

	if db.Error != nil {
		entry.Error("Unable to query record")
	}
	entry.Trace("Unable to query record")
	return db.Error
}

func (repo *CommonRepositoryDB) AddPagination(db *gorm.DB, logFields log.Fields, query PaginationQuery) *gorm.DB {
	if query.Limit > 0 {
		logFields["offset"] = query.Offset
		logFields["limit"] = query.Limit

		db = db.Offset(query.Offset).Limit(query.Limit)
	}
	return db
}
