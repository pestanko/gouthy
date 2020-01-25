package services

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type CommonService struct {
	DB gorm.DB
	Name string
}

func NewCommonService(db gorm.DB, name string) CommonService {
	return CommonService{DB: db, Name: name}
}

func (service *CommonService) Create(instance interface{}) error {
	result := service.DB.Create(instance)
	if result.Error != nil {
		log.Error("Create {} failed: ", service.Name, result.Error)
		return result.Error
	}
	log.Info("Created {}", service.Name, instance)
	return nil
}

func (service *CommonService) Update(instance interface{}) error {
	result := service.DB.Update(instance)
	if result.Error != nil {
		log.Error("Update {} failed: ", service.Name, result.Error)
		return result.Error
	}
	log.Info("Updated {}", service.Name, instance)
	return nil
}

func (service *CommonService) Delete(instance interface{}) error {
	result := service.DB.Delete(instance)
	if result.Error != nil {
		log.Error("Delete {} failed: ", service.Name, result.Error)
		return result.Error
	}
	log.Info("Deleted {}", service.Name, instance)
	return nil
}





