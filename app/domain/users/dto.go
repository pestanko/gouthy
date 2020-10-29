package users

import (
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

// DTO

type UserDTO struct {
	ID        uuid.UUID `json:"id" yaml:"id"`
	Username  string    `json:"username" yaml:"username"`
	Name      string    `json:"name" yaml:"name"`
	Email     string    `json:"email" yaml:"email"`
	State     string    `json:"state" yaml:"state"`
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at"`
}

type CreateDTO struct {
	Username string `json:"username" yaml:"username"`
	Name     string `json:"name" yaml:"name"`
	Email    string `json:"email" yaml:"email"`
	Password string `json:"password" yaml:"password"`
}

func (dto *CreateDTO) ToEntity() User {
	return User{
		ID:       uuid.Nil,
		Username: dto.Username,
		Name:     dto.Name,
		Email:    dto.Email,
	}
}

func (dto *CreateDTO) LogFields() log.Fields {
	return log.Fields{
		"user_username": dto.Username,
		"user_email":    dto.Email,
	}
}

type UpdateDTO struct {
	Username string `json:"username" yaml:"username"`
	Name     string `json:"name" yaml:"name"`
	Email    string `json:"email" yaml:"email"`
}

func (dto *UpdateDTO) ToEntity() User {
	return User{
		Username: dto.Username,
		Name:     dto.Name,
		Email:    dto.Email,
	}
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
