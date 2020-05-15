package entities

import (
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

type Facade interface {
	List() ([]ListEntity, error)
	Get(userId uuid.UUID) (*Entity, error)
}

type FacadeImpl struct {
	entities Repository
	secrets  SecretsRepository
}

func (s *FacadeImpl) List() ([]ListEntity, error) {
	panic("implement me")
}

func (s *FacadeImpl) Get(userId uuid.UUID) (*Entity, error) {
	panic("implement me")
}

func NewEntitiesFacade(entities Repository, secrets SecretsRepository) Facade {
	return &FacadeImpl{
		entities: entities,
		secrets:  secrets,
	}
}
