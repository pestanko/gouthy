package services

import (
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/models"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type UsersService struct {
	DB     *gorm.DB
	common CommonService
}

func NewUsersService(db *gorm.DB) UsersService {
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

func (service *UsersService) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	result := service.DB.Find(&user).Where("id = ?", id)
	if result.Error != nil {
		log.WithFields(log.Fields{
			"id": id,
		}).WithError(result.Error).Error("Find Failed", id)
		return nil, result.Error
	}
	return &user, nil
}
