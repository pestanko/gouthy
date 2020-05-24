package users

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
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	UserID    uuid.UUID  `gorm:"type:varchar" json:"user_id"`
	Value     string     `gorm:"type:varchar" json:"value"`
	Name      string     `gorm:"type:varchar" json:"name"`
	CreatedAt time.Time  `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:timestamp" json:"updated_at"`
	ExpiresAt *time.Time `gorm:"type:timestamp" json:"expires_at"`
}

func (Secret) TableName() string {
	return "userSecrets"
}

func (s Secret) IsExpired() bool {
	return s.ExpiresAt != nil && s.ExpiresAt.Before(time.Now())
}

type SecretsRepository interface {
	Create(ctx context.Context, secret *Secret) error
	Update(ctx context.Context, secret *Secret) error
	Delete(ctx context.Context, secret *Secret) error
	FindByID(ctx context.Context, userId uuid.UUID, id uuid.UUID) (*Secret, error)
	List(ctx context.Context, userId uuid.UUID) ([]Secret, error)
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

func (r *secretsRepositoryDB) FindByID(ctx context.Context, userId uuid.UUID, id uuid.UUID) (*Secret, error) {
	var secret Secret
	result := r.DB.Where("id = ? AND user_id = ?", id, userId).Find(&secret)
	if result.Error != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"id":      id,
			"user_id": userId,
		}).WithError(result.Error).Error("Find Failed")

		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, nil
		}

		return nil, result.Error
	}
	return &secret, nil
}

func (r *secretsRepositoryDB) List(ctx context.Context, userId uuid.UUID) ([]Secret, error) {
	var secrets []Secret
	r.DB.Where("user_id = ?", userId).Find(&secrets)
	return secrets, r.DB.Error
}

func NewSecretsRepositoryDB(DB *gorm.DB) SecretsRepository {
	return &secretsRepositoryDB{DB: DB, common: repositories.NewCommonRepositoryDB(DB, "userSecrets")}
}
