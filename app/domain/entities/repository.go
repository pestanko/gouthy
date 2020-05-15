package entities

import (
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/shared/repositories"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type Repository interface {
	Create(entity *Entity) error
	Update(entity *Entity) error
	Delete(entity *Entity) error
	List() ([]Entity, error)
	FindByID(id uuid.UUID) (*Entity, error)
}

type RepositoryDB struct {
	db     *gorm.DB
	common repositories.CommonRepository
}

func NewEntitiesRepositoryDB(db *gorm.DB) Repository {
	return &RepositoryDB{db: db, common: repositories.NewCommonRepository(db, "Entity")}
}

func (service *RepositoryDB) Create(entity *Entity) error {
	return service.common.Create(entity)
}

func (service *RepositoryDB) Update(entity *Entity) error {
	return service.common.Update(entity)
}

func (service *RepositoryDB) Delete(entity *Entity) error {
	return service.common.Delete(entity)
}

func (service *RepositoryDB) List() ([]Entity, error) {
	var entities []Entity
	result := service.db.Find(&entities)
	if result.Error != nil {
		log.Error("List failed: {}", result.Error)
		return nil, result.Error
	}
	return entities, nil
}

func (service *RepositoryDB) FindByID(id uuid.UUID) (*Entity, error) {
	var entity Entity
	result := service.db.Find(&entity).Where("id = ?", id)
	if result.Error != nil {
		log.WithFields(log.Fields{
			"id": id,
		}).WithError(result.Error).Error("Find Failed", id)
		return nil, result.Error
	}
	return &entity, nil
}