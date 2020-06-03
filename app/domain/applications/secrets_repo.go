package applications

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/shared/repositories"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

type Secret struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	ApplicationId uuid.UUID  `gorm:"type:varchar" json:"application_id"`
	Value         string     `gorm:"type:varchar" json:"value"`
	CreatedAt     time.Time  `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"type:timestamp" json:"updated_at"`
	ExpiresAt     *time.Time `gorm:"type:timestamp" json:"expires_at"`
}

type SecretsRepository interface {
	Create(ctx context.Context, secret *Secret) error
	Update(ctx context.Context, secret *Secret) error
	Delete(ctx context.Context, secret *Secret) error
	FindByID(ctx context.Context, appId uuid.UUID, id uuid.UUID) (*Secret, error)
	List(ctx context.Context, appId uuid.UUID) ([]Secret, error)
}

type secretsRepositoryDB struct {
	DB     *gorm.DB
	common repositories.CommonRepositoryDB
}

func (r *secretsRepositoryDB) Create(ctx context.Context, secret *Secret) error {
	return r.common.Create(ctx, secret)
}

func (r *secretsRepositoryDB) Update(ctx context.Context, secret *Secret) error {
	return r.common.Update(ctx, secret)
}

func (r *secretsRepositoryDB) Delete(ctx context.Context, secret *Secret) error {
	return r.common.Delete(ctx, secret)
}

func (r *secretsRepositoryDB) FindByID(ctx context.Context, appId uuid.UUID, id uuid.UUID) (*Secret, error) {
	var secret Secret
	result := r.DB.Where("id = ? AND application_id = ?", id, appId).Find(&secret)
	if result.Error != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"id":             id,
			"application_id": appId,
		}).WithError(result.Error).Error("Find Failed")

		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, nil
		}

		return nil, result.Error
	}
	return &secret, nil
}

func (r *secretsRepositoryDB) List(ctx context.Context, appId uuid.UUID) ([]Secret, error) {
	var secrets []Secret
	r.DB.Where("application_id = ?", appId).Find(&secrets)
	return secrets, r.DB.Error
}

func NewSecretsRepositoryDB(DB *gorm.DB) SecretsRepository {
	return &secretsRepositoryDB{DB: DB, common: repositories.NewCommonRepositoryDB(DB, "ApplicationSecrets")}
}

func (Secret) TableName() string {
	return "ApplicationSecrets"
}
