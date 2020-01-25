package services

import (
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/models"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type AccountsService struct {
	DB gorm.DB
	common CommonService
}

type SecretsService struct {
	DB gorm.DB
	common CommonService
}

func NewAccountsService(db gorm.DB) AccountsService {
	return AccountsService{DB: db, common: NewCommonService(db, "Account") }
}

func NewSecretsService(db gorm.DB) SecretsService {
	return SecretsService{DB: db, common: NewCommonService(db, "Secret")}
}

func (service *AccountsService) Create(account *models.Account) error {
	return service.common.Create(account)
}

func (service *AccountsService) Update(account *models.Account) error {
	return service.common.Update(account)

}

func (service *AccountsService) Delete(account *models.Account) error {
	return service.common.Delete(account)
}

func (service *AccountsService) List() ([]models.Account, error) {
	var accounts []models.Account
	result := service.DB.Find(&accounts)
	if result.Error != nil {
		log.Error("List failed: {}", result.Error)
		return nil, result.Error
	}
	return accounts, nil
}

func (service *AccountsService) FindByID(id uuid.UUID) (*models.Account, error) {
	var account models.Account
	result := service.DB.Find(&account).Where("id = ?", id)
	if result.Error != nil {
		log.Error("Find for id {} failed: {}", id, result.Error)
		return nil, result.Error
	}
	return &account, nil
}

func (service *AccountsService) FindByEntityId(id uuid.UUID) (*models.Account, error) {
	var account models.Account
	result := service.DB.Find(&account).Where("entity_id = ?", id)
	if result.Error != nil {
		log.Error("Find for entity_id {} failed: {}", id, result.Error)
		return nil, result.Error
	}
	return &account, nil
}


