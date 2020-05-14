package entities

import (
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/shared/repositories"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Secret struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	EntityId  uuid.UUID `gorm:"type:uuid" json:"entity_id"`
	Name      string    `gorm:"varchar" json:"name"`
	Value     string    `gorm:"varchar" json:"-"`
	CreatedAt time.Time `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp" json:"updated_at"`
	ExpiresAt time.Time `gorm:"type:timestamp" json:"expires_at"`
}

type SecretsRepository struct {
	DB     *gorm.DB
	common repositories.CommonRepository
}

func NewSecretsRepository(db *gorm.DB) SecretsRepository {
	return SecretsRepository{DB: db, common: repositories.NewCommonService(db, "Secret")}
}
