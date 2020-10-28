package apps

import (
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

const DefaultApplicationClientId = "default"
const ConsoleApplicationClientId = "admin_console"

// DTO

type AppDTO struct {
	baseApplicationDTO
	Description     string     `json:"description"`
	ClientId        string     `json:"client_id"`
	State           string     `json:"state"`
	Type            string     `json:"type"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	RedirectUris    []string   `json:"redirect_uris"`
	AvailableScopes []string   `json:"available_scopes"`
}

type ListApplicationDTO struct {
	baseApplicationDTO
}

type CreateDTO struct {
	Codename    string `json:"codename"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ClientId    string `json:"client_id"`
	Type        string `json:"type"`
}

func (d *CreateDTO) LogFields() log.Fields {
	return log.Fields{
		"codename":  d.Codename,
		"client_id": d.ClientId,
	}
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
