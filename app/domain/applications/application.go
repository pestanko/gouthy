package applications

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Application struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt   time.Time  `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"type:timestamp" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"type:timestamp" json:"deleted_at"`
	Codename    string     `gorm:"varchar" json:"codename"`
	Name        string     `gorm:"varchar" json:"name"`
	Type        string     `gorm:"varchar" json:"type"`
	Description string     `gorm:"varchar" json:"description"`
	ClientId    string     `gorm:"varchar" json:"client_id"`
	State       string     `gorm:"varchar" json:"state"`
}

// DTO
type baseApplicationDTO struct {
	Codename string    `json:"username"`
	Name     string    `json:"name"`
	ID       uuid.UUID `json:"id"`
}

type ListApplicationDTO struct {
	baseApplicationDTO
}

type CreateDTO struct {
	Codename    string `gorm:"varchar" json:"codename"`
	Name        string `gorm:"varchar" json:"name"`
	Description string `gorm:"varchar" json:"description"`
	ClientId    string `gorm:"varchar" json:"client_id"`
	Type        string `gorm:"varchar" json:"type"`
}

type UpdateDTO struct {
	Codename    string `gorm:"varchar" json:"codename"`
	Name        string `gorm:"varchar" json:"name"`
	Description string `gorm:"varchar" json:"description"`
}

type ApplicationDTO struct {
	baseApplicationDTO
	Description string     `json:"description"`
	ClientId    string     `json:"client_id"`
	State       string     `json:"state"`
	Type        string     `json:"type"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
