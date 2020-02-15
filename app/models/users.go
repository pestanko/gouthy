package models

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

type User struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	AccountId uuid.UUID `gorm:"type:uuid" json:"account_id"`
	CreatedAt time.Time `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp" json:"updated_at"`
	Username  string    `gorm:"type:varchar" json:"username"`
	Password  string    `gorm:"type:varchar" json:"password_hash"`
	Name      string    `gorm:"type:varchar" json:"name"`
	Email     string    `gorm:"type:varchar" json:"email"`
}

type ForgotPasswordCode struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Code      string    `gorm:"type:varchar" json:"code"`
	UserId    uuid.UUID `gorm:"type:uuid" json:"user_id"`
	CreatedAt time.Time `gorm:"type:timestamp" json:"created_at"`
	UsedAt    time.Time `gorm:"type:timestamp" json:"used_at"`
}
