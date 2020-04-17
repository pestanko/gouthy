package services

import (
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/models"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type EntitiesService struct {
	DB     *gorm.DB
	common CommonService
}

type SecretsService struct {
	DB     *gorm.DB
	common CommonService
}

func NewAccountsService(db *gorm.DB) EntitiesService {
	return EntitiesService{DB: db, common: NewCommonService(db, "Account")}
}

func NewSecretsService(db *gorm.DB) SecretsService {
	return SecretsService{DB: db, common: NewCommonService(db, "Secret")}
}

func (service *EntitiesService) Create(entity *models.Entity) error {
	return service.common.Create(entity)
}

func (service *EntitiesService) Update(entity *models.Entity) error {
	return service.common.Update(entity)

}

func (service *EntitiesService) Delete(entity *models.Entity) error {
	return service.common.Delete(entity)
}

func (service *EntitiesService) List() ([]models.Entity, error) {
	var entities []models.Entity
	result := service.DB.Find(&entities)
	if result.Error != nil {
		log.Error("List failed: {}", result.Error)
		return nil, result.Error
	}
	return entities, nil
}

func (service *EntitiesService) FindByID(id uuid.UUID) (*models.Entity, error) {
	var entity models.Entity
	result := service.DB.Find(&entity).Where("id = ?", id)
	if result.Error != nil {
		log.WithFields(log.Fields{
			"id": id,
		}).WithError(result.Error).Error("Find Failed", id)
		return nil, result.Error
	}
	return &entity, nil
}