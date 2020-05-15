package users

import (
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/shared/repositories"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type Repository struct {
	DB     *gorm.DB
	common repositories.CommonRepositoryDB
}

func NewUsersRepositoryDB(db *gorm.DB) Repository {
	return Repository{DB: db, common: repositories.NewCommonRepositoryDB(db, "User")}
}

func (r *Repository) Create(user *User) error {
	return r.common.Create(user)
}

func (r *Repository) Update(user *User) error {
	return r.common.Update(user)
}

func (r *Repository) Delete(user *User) error {
	return r.common.Delete(user)
}

func (r *Repository) FindByID(id uuid.UUID) (*User, error) {
	var user User
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

func (r *Repository) List() ([]User, error) {
	var users []User
	r.DB.Find(&users)

	return users, r.DB.Error
}

func (r *Repository) FindByUsername(username string) (*User, error) {
	var user User
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
