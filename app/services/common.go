package services

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type CommonService struct {
	DB   *gorm.DB
	Name string
}

func NewCommonService(db *gorm.DB, name string) CommonService {
	return CommonService{DB: db, Name: name}
}

func (service *CommonService) Create(instance interface{}) error {
	result := service.DB.Create(instance)
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

func (service *CommonService) Update(instance interface{}) error {
	result := service.DB.Update(instance)
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

func (service *CommonService) Delete(instance interface{}) error {
	result := service.DB.Delete(instance)
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
