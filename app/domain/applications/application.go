package applications

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type ApplicationModel struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;"`
	CreatedAt   time.Time  `gorm:"type:timestamp"`
	UpdatedAt   time.Time  `gorm:"type:timestamp"`
	DeletedAt   *time.Time `gorm:"type:timestamp"`
	Codename    string     `gorm:"varchar"`
	Name        string     `gorm:"varchar"`
	Type        string     `gorm:"varchar"`
	Description string     `gorm:"varchar"`
	ClientId    string     `gorm:"varchar"`
	State       string     `gorm:"varchar"`
}

func (ApplicationModel) TableName() string {
	return "Applications"
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
