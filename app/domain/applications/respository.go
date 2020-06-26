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

type applicationQuery struct {
	repositories.PaginationQuery

	Id       uuid.UUID
	Codename string
	ClientId string
}

type applicationModel struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;"`
	CreatedAt   time.Time  `gorm:"type:timestamp"`
	UpdatedAt   time.Time  `gorm:"type:timestamp"`
	DeletedAt   *time.Time `gorm:"type:timestamp"`
	Codename    string     `gorm:"varchar"`
	Name        string     `gorm:"varchar"`
	Type        string     `gorm:"varchar"`
	Description string     `gorm:"varchar"`
	ClientId    string     `gorm:"varchar"`
	State       string     `gorm:"varchar"`
}

func (applicationModel) TableName() string {
	return "Applications"
}

type Repository interface {
	Create(ctx context.Context, app *applicationModel) error
	Update(ctx context.Context, app *applicationModel) error
	Delete(ctx context.Context, app *applicationModel) error
	Query(ctx context.Context, query applicationQuery) ([]applicationModel, error)
	QueryOne(ctx context.Context, query applicationQuery) (*applicationModel, error)

	List(ctx context.Context) ([]applicationModel, error)
}

type repositoryDB struct {
	DB     *gorm.DB
	common repositories.CommonRepositoryDB
}

func (r *repositoryDB) Query(ctx context.Context, query applicationQuery) (result []applicationModel, err error) {
	db, entry := r.internalQueryBuilder(ctx, query)
	return result, r.common.ProcessQuery(db, &result, entry)
}

func (r *repositoryDB) QueryOne(ctx context.Context, query applicationQuery) (*applicationModel, error) {
	var result applicationModel
	db, entry := r.internalQueryBuilder(ctx, query)
	one, err := r.common.ProcessQueryOne(db, &result, entry)
	if one == nil {
		return nil, err
	}
	return one.(*applicationModel), err
}

func (r *repositoryDB) Create(ctx context.Context, user *applicationModel) error {
	return r.common.Create(ctx, user)
}

func (r *repositoryDB) Update(ctx context.Context, user *applicationModel) error {
	return r.common.Update(ctx, user)
}

func (r *repositoryDB) Delete(ctx context.Context, user *applicationModel) error {
	return r.common.Delete(ctx, user)
}

func (r *repositoryDB) FindByID(ctx context.Context, id uuid.UUID) (*applicationModel, error) {
	var app applicationModel
	result := r.DB.Where("id = ?", id).Find(&app)
	if result.Error != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
			"id": id,
		}).WithError(result.Error).Error("Find Failed")

		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, nil
		}

		return nil, result.Error
	}
	return &app, nil
}

func (r *repositoryDB) List(ctx context.Context) (result []applicationModel, err error) {
	r.DB.Find(&result)
	if r.DB.Error != nil {
		shared.GetLogger(ctx).WithFields(log.Fields{
		}).WithError(err).Error("List applications failed")
	}
	return result, r.DB.Error
}

func (r *repositoryDB) internalQueryBuilder(ctx context.Context, query applicationQuery) (*gorm.DB, *log.Entry) {
	db := r.DB
	logFields := log.Fields{
		"model": "user",
	}
	if query.Id != uuid.Nil {
		db = db.Where("id = ?", query.Id)
		logFields["id"] = query.Id
	}
	if query.ClientId != "" {
		db = db.Where("client_id = ?", query.ClientId)
		logFields["client_id"] = query.ClientId
	}

	if query.Codename != "" {
		db = db.Where("codename = ?", query.Codename)
		logFields["username"] = query.Codename
	}

	db = r.common.AddPagination(db, logFields, query.PaginationQuery)

	return db, shared.GetLogger(ctx).WithFields(logFields)
}

func NweApplicationsRepositoryDB(db *gorm.DB) Repository {
	return &repositoryDB{
		DB:     db,
		common: repositories.NewCommonRepositoryDB(db, "Applications"),
	}
}
