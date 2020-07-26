package users

import (
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

// DTO

type UserDTO struct {
	baseUserDTO
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type baseUserDTO struct {
	Username string    `json:"Username"`
	Name     string    `json:"name"`
	Email    string    `json:"Email"`
	ID       uuid.UUID `json:"Id"`
}

type CreateDTO struct {
	baseUserDTO
	Password string `json:"password"`
}

func (d *CreateDTO) LogFields() log.Fields {
	return log.Fields{
		"username": d.Username,
		"user_id": d.ID,
	}
}

type UpdateDTO struct {
	Username string `json:"Username"`
	Name     string `json:"name"`
	Email    string `json:"Email"`
}

type UpdatePasswordDTO struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ListUserDTO struct {
	baseUserDTO
}
