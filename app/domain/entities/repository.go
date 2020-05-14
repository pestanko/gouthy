package entities

import (
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/shared/repositories"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type Repository struct {
	DB     *gorm.DB
	common repositories.CommonRepository
}

func NewEntitiesRepository(db *gorm.DB) Repository {
	return Repository{DB: db, common: repositories.NewCommonService(db, "Entity")}
}

func (service *Repository) Create(entity *Entity) error {
	return service.common.Create(entity)
}

func (service *Repository) Update(entity *Entity) error {
	return service.common.Update(entity)
}

func (service *Repository) Delete(entity *Entity) error {
	return service.common.Delete(entity)
}

func (service *Repository) List() ([]Entity, error) {
	var entities []Entity
	result := service.DB.Find(&entities)
	if result.Error != nil {
		log.Error("List failed: {}", result.Error)
		return nil, result.Error
	}
	return entities, nil
}

func (service *Repository) FindByID(id uuid.UUID) (*Entity, error) {
	var entity Entity
	result := service.DB.Find(&entity).Where("id = ?", id)
	if result.Error != nil {
		log.WithFields(log.Fields{
			"id": id,
		}).WithError(result.Error).Error("Find Failed", id)
		return nil, result.Error
	}
	return &entity, nil
}