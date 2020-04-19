package models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Entity struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Type      string    `gorm:"type:varchar" json:"entity_type"`
	State     string    `gorm:"type:varchar" json:"entity_state"`
	CreatedAt time.Time `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp" json:"updated_at"`
	DeletedAt time.Time `gorm:"type:timestamp" json:"deleted_at"`
}

type Secret struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	EntityId  uuid.UUID `gorm:"type:uuid" json:"entity_id"`
	Name      string    `gorm:"varchar" json:"name"`
	Value     string    `gorm:"varchar" json:"-"`
	CreatedAt time.Time `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp" json:"updated_at"`
	ExpiresAt time.Time `gorm:"type:timestamp" json:"expires_at"`
}

type AutomaticSecurityCodes struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Code      string    `gorm:"type:varchar" json:"code"`
	EntityId  uuid.UUID `gorm:"type:uuid" json:"entity_id"`
	CreatedAt time.Time `gorm:"type:timestamp" json:"created_at"`
	UsedAt    time.Time `gorm:"type:timestamp" json:"used_at"`
}

func NewEntity() *Entity {
	return &Entity{State: "created"}
}

// Entities is not Entitys
func (Entity) TableName() string {
	return "entities"
}
