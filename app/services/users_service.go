package services

import (
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/models"
	"github.com/pestanko/gouthy/app/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
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
	hash, err := utils.HashString(password)
	if err != nil {
		log.Error("Unable to hash a password: {}", err)
		return err
	}

	user.Password = hash

	if err := service.Update(user); err != nil {
		return err
	}

	return nil
}

func (service *UsersService) CheckPassword(user *models.User, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
}

