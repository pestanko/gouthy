package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/models"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type UsersRepository struct {
	DB     *gorm.DB
	common CommonRepository
}

func NewUsersRepository(db *gorm.DB) UsersRepository {
	return UsersRepository{DB: db, common: NewCommonService(db, "User")}
}

func (r *UsersRepository) Create(user *models.User) error {
	return r.common.Create(user)
}

func (r *UsersRepository) Update(user *models.User) error {
	return r.common.Update(user)
}

func (r *UsersRepository) Delete(user *models.User) error {
	return r.common.Delete(user)
}

func (r *UsersRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	result := r.DB.Where("id = ?", id).Find(&user)
	if result.Error != nil {
		log.WithFields(log.Fields{
			"id": id,
		}).WithError(result.Error).Error("Find Failed")

		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, nil
		}

		return nil, result.Error
	}
	return &user, nil
}

func (r *UsersRepository) List() ([]models.User, error) {
	var users []models.User
	r.DB.Find(&users)

	return users, r.DB.Error
}

func (r *UsersRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.DB.Where("username = ?", username).Find(&user)
	if result.Error != nil {
		log.WithFields(log.Fields{
			"username": username,
		}).WithError(result.Error).Error("Find Failed")

		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, nil
		}

		return nil, result.Error
	}
	return &user, nil
}
