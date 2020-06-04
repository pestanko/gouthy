package users

import (
	"github.com/pestanko/gouthy/app/shared/utils"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt time.Time  `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:timestamp" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:timestamp" json:"deleted_at"`
	Username  string     `gorm:"type:varchar" json:"username"`
	Password  string     `gorm:"type:varchar" json:"-"`
	Name      string     `gorm:"type:varchar" json:"name"`
	Email     string     `gorm:"type:varchar" json:"email"`
}

func (UserModel) TableName() string {
	return "Users"
}

func (user *UserModel) SetPassword(password string) error {
	hash, err := utils.HashString(password)
	if err != nil {
		return err
	}

	user.Password = hash
	return nil
}

func (user *UserModel) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
}

// DTO
type baseUserDTO struct {
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	ID       uuid.UUID `json:"id"`
}

type CreateDTO struct {
	baseUserDTO
	Password string `json:"password"`
}

type UpdateDTO struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type UpdatePasswordDTO struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ListUserDTO struct {
	baseUserDTO
}

type User struct {
	baseUserDTO
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (user *UserModel) ToEntity() *User {
	return &User{
		baseUserDTO: *convertModelToUserBase(user),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
