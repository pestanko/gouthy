package models

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Account struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	EntityId    uuid.UUID `gorm:"type:uuid" json:"entity_id"`
	AccountType  string    `gorm:"type:varchar" json:"account_type"`
	AccountState string    `gorm:"type:varchar" json:"account_state"`
	CreatedAt    time.Time `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp" json:"updated_at"`
}

type Secret struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	AccountId uuid.UUID `gorm:"type:uuid" json:"account_id"`
	Name      string    `gorm:"varchar" json:"name"`
	Value     string    `gorm:"varchar" json:"value"`
	CreatedAt time.Time `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp" json:"updated_at"`
	ExpiresAt time.Time `gorm:"type:timestamp" json:"expires_at"`
}
