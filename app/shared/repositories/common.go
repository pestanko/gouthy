package repositories

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type CommonRepository struct {
	db   *gorm.DB
	Name string
}

func NewCommonService(db *gorm.DB, name string) CommonRepository {
	return CommonRepository{db: db, Name: name}
}

func (service *CommonRepository) Create(instance interface{}) error {
	result := service.db.Create(instance)
	if result.Error != nil {
		log.WithFields(log.Fields{
			"service":  service.Name,
			"action":   "create",
			"instance": instance,
		}).WithError(result.Error).Error("Create failed")
		return result.Error
	}
	log.WithFields(log.Fields{
		"service":  service.Name,
		"action":   "create",
		"instance": instance,
	}).Info("Created instance")
	return nil
}

func (service *CommonRepository) Update(instance interface{}) error {
	result := service.db.Update(instance)
	if result.Error != nil {
		log.WithFields(log.Fields{
			"service":  service.Name,
			"action":   "update",
			"instance": instance,
		}).WithError(result.Error).Error("Update failed")
		return result.Error
	}
	log.WithFields(log.Fields{
		"service":  service.Name,
		"action":   "update",
		"instance": instance,
	}).Info("Updated instance")
	return nil
}

func (service *CommonRepository) Delete(instance interface{}) error {
	result := service.db.Delete(instance)
	if result.Error != nil {
		log.WithFields(log.Fields{
			"service":  service.Name,
			"action":   "delete",
			"instance": instance,
		}).WithError(result.Error).Error("Delete failed")
		return result.Error
	}
	log.WithFields(log.Fields{
		"service":  service.Name,
		"action":   "delete",
		"instance": instance,
	}).Info("Deleted instance")
	return nil
}
