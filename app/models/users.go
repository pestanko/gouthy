package models

import (
	"github.com/pestanko/gouthy/app/utils"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt time.Time `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp" json:"updated_at"`
	Username  string    `gorm:"type:varchar" json:"username"`
	Password  string    `gorm:"type:varchar" json:"-"`
	Name      string    `gorm:"type:varchar" json:"name"`
	Email     string    `gorm:"type:varchar" json:"email"`
}

func (user *User) SetPassword(password string) error {
	hash, err := utils.HashString(password)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"entity_id": user.ID,
			"username":  user.Username,
		}).Error("Unable to hash a password")
		return err
	}

	user.Password = hash
	return nil
}

func (user *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
}
