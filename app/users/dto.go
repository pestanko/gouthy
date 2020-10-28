package users

import (
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

// DTO

type UserDTO struct {
	baseUserDTO
	ID        uuid.UUID `json:"id" yaml:"id"`
	State     string    `json:"state" yaml:"state"`
	CreatedAt time.Time `json:"created_at" yaml:"state"`
	UpdatedAt time.Time `json:"updated_at" yaml:"state"`
}

type baseUserDTO struct {
	Username string `json:"username" yaml:"state"`
	Name     string `json:"name" yaml:"state"`
	Email    string `json:"email" yaml:"email"`
}

type CreateDTO struct {
	baseUserDTO
	Password string `json:"password"`
}

func (d *CreateDTO) LogFields() log.Fields {
	return log.Fields{
		"user_username": d.Username,
		"user_email":    d.Email,
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
	ID       uuid.UUID `json:"id" yaml:"id"`
	Username string    `json:"username" yaml:"state"`
	Email    string    `json:"email" yaml:"email"`
}
