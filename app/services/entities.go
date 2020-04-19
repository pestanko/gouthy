package services

import (
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/models"
	"github.com/pestanko/gouthy/app/repositories"
	uuid "github.com/satori/go.uuid"
	"time"
)

type EntityBase struct {
	ID   uuid.UUID `json:"id"`
	Type string    `json:"type"`
}

type ListEntity struct {
	EntityBase
}

type ListSecrets struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type ListEntitySecrets struct {
	EntityBase
	Secrets []ListSecrets
}

type EntityDTO struct {
	EntityBase
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type EntitiesService interface {
	List() ([]ListEntity, error)
	Get(userId uuid.UUID) (*models.Entity, error)
}

type EntitiesServiceImpl struct {
	users    repositories.UsersRepository
	entities repositories.EntitiesRepository
	secrets  repositories.SecretsService
}

func (s EntitiesServiceImpl) List() ([]ListEntity, error) {
	panic("implement me")
}

func (s EntitiesServiceImpl) Get(userId uuid.UUID) (*models.Entity, error) {
	panic("implement me")
}

func NewEntitiesService(db *gorm.DB) EntitiesService {
	return &EntitiesServiceImpl{
		users:    repositories.NewUsersRepository(db),
		entities: repositories.NewEntitiesRepository(db),
		secrets:  repositories.NewSecretsRepository(db),
	}
}
