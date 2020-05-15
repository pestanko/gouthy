package entities

import (
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/shared/repositories"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Secret struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	EntityId  uuid.UUID  `gorm:"type:uuid" json:"entity_id"`
	Name      string     `gorm:"varchar" json:"name"`
	Value     string     `gorm:"varchar" json:"-"`
	CreatedAt time.Time  `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:timestamp" json:"updated_at"`
	ExpiresAt *time.Time `gorm:"type:timestamp" json:"expires_at"`
}

func (s Secret) IsExpired() bool {
	now := time.Now()

	return s.ExpiresAt != nil && now.After(*s.ExpiresAt)
}

type SecretsRepository interface {
	FindSecretsForEntity(id uuid.UUID) ([]Secret, error)
}

type SecretsRepositoryDB struct {
	db     *gorm.DB
	common repositories.CommonRepositoryDB
}

func (repo *SecretsRepositoryDB) FindSecretsForEntity(id uuid.UUID) ([]Secret, error) {
	panic("implement me")
}

func NewSecretsRepositoryDB(db *gorm.DB) SecretsRepository {
	return &SecretsRepositoryDB{db: db, common: repositories.NewCommonRepositoryDB(db, "Secret")}
}
