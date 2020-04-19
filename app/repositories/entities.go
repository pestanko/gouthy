package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/models"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type EntitiesRepository struct {
	DB     *gorm.DB
	common CommonRepository
}

type SecretsService struct {
	DB     *gorm.DB
	common CommonRepository
}

func NewEntitiesRepository(db *gorm.DB) EntitiesRepository {
	return EntitiesRepository{DB: db, common: NewCommonService(db, "Entity")}
}

func NewSecretsRepository(db *gorm.DB) SecretsService {
	return SecretsService{DB: db, common: NewCommonService(db, "Secret")}
}

func (service *EntitiesRepository) Create(entity *models.Entity) error {
	return service.common.Create(entity)
}

func (service *EntitiesRepository) Update(entity *models.Entity) error {
	return service.common.Update(entity)

}

func (service *EntitiesRepository) Delete(entity *models.Entity) error {
	return service.common.Delete(entity)
}

func (service *EntitiesRepository) List() ([]models.Entity, error) {
	var entities []models.Entity
	result := service.DB.Find(&entities)
	if result.Error != nil {
		log.Error("List failed: {}", result.Error)
		return nil, result.Error
	}
	return entities, nil
}

func (service *EntitiesRepository) FindByID(id uuid.UUID) (*models.Entity, error) {
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