package services

import (
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/models"
	"golang.org/x/crypto/bcrypt"
	log "github.com/sirupsen/logrus"
)

type UsersService struct {
	DB     gorm.DB
	common CommonService
}

func NewUsersService(db gorm.DB) UsersService {
	return UsersService{DB: db, common: NewCommonService(db, "User")}
}

func (service *UsersService) Create(user *models.User) error {
	return service.common.Create(user)
}

func (service *UsersService) Update(user *models.User) error {
	return service.common.Update(user)
}

func (service *UsersService) Delete(user *models.User) error {
	return service.common.Delete(user)
}

func (service *UsersService) SetPassword(user *models.User, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		log.Error("Unable to hash a password: {}", err)
		return err
	}
	user.Password = string(hash)

	if err := service.Update(user); err != nil {
		return err
	}

	return nil
}

func (service *UsersService) CheckPassword(user *models.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err != nil
}

