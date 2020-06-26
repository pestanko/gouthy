package applications

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

// DTO

type Application struct {
	baseApplicationDTO
	Description string     `json:"description"`
	ClientId    string     `json:"client_id"`
	State       string     `json:"state"`
	Type        string     `json:"type"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
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

type baseApplicationDTO struct {
	Codename string    `json:"username"`
	Name     string    `json:"name"`
	ID       uuid.UUID `json:"id"`
}
