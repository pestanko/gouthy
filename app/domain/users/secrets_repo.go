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

type SecretQuery struct {
	repositories.PaginationQuery
	Id     uuid.UUID
	UserId uuid.UUID
}

type SecretModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;"`
	UserID    uuid.UUID  `gorm:"type:varchar"`
	Value     string     `gorm:"type:varchar"`
	Name      string     `gorm:"type:varchar"`
	CreatedAt time.Time  `gorm:"type:timestamp"`
	UpdatedAt time.Time  `gorm:"type:timestamp"`
	ExpiresAt *time.Time `gorm:"type:timestamp"`
}

func (SecretModel) TableName() string {
	return "UserSecrets"
}

func (s SecretModel) IsExpired() bool {
	return s.ExpiresAt != nil && s.ExpiresAt.Before(time.Now())
}

type SecretsRepository interface {
	Create(ctx context.Context, secret *SecretModel) error
	Update(ctx context.Context, secret *SecretModel) error
	Delete(ctx context.Context, secret *SecretModel) error
	Query(ctx context.Context, query SecretQuery) ([]SecretModel, error)
	QueryOne(ctx context.Context, query SecretQuery) (*SecretModel, error)
}

type secretsRepositoryDB struct {
	DB     *gorm.DB
	common repositories.CommonRepositoryDB
}

func (r *secretsRepositoryDB) QueryOne(ctx context.Context, query SecretQuery) (*SecretModel, error) {
	var result SecretModel
	db, entry := r.internalQueryBuilder(ctx, query)
	one, err := r.common.ProcessQueryOne(db, &result, entry)
	if one == nil {
		return nil, err
	}
	return  one.(*SecretModel), err
}

func (r *secretsRepositoryDB) Query(ctx context.Context, query SecretQuery) (result []SecretModel, err error) {
	db, entry := r.internalQueryBuilder(ctx, query)
	return result, r.common.ProcessQuery(db, &result, entry)
}

func (r *secretsRepositoryDB) Create(ctx context.Context, secret *SecretModel) error {
	return r.common.Create(ctx, secret)
}

func (r *secretsRepositoryDB) Update(ctx context.Context, secret *SecretModel) error {
	return r.common.Update(ctx, secret)
}

func (r *secretsRepositoryDB) Delete(ctx context.Context, secret *SecretModel) error {
	return r.common.Delete(ctx, secret)
}

func (r *secretsRepositoryDB) internalQueryBuilder(ctx context.Context, query SecretQuery) (*gorm.DB, *log.Entry) {
	db := r.DB
	logFields := log.Fields{
		"record": "user_secret",
	}

	if query.Id != uuid.Nil {
		db = db.Where("id = ?", query.Id)
		logFields["id"] = query.Id
	}

	if query.UserId != uuid.Nil {
		db = db.Where("user_id = ?", query.UserId)
		logFields["user_id"] = query.UserId
	}

	db = r.common.AddPagination(db, logFields, query.PaginationQuery)

	return db, shared.GetLogger(ctx).WithFields(logFields)
}

func NewSecretsRepositoryDB(DB *gorm.DB) SecretsRepository {
	return &secretsRepositoryDB{DB: DB, common: repositories.NewCommonRepositoryDB(DB, "userSecrets")}
}
