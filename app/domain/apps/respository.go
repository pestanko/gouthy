package apps

import (
	"context"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/shared/repos"
	"github.com/pestanko/gouthy/app/shared/utils"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"time"
)

type FindQuery struct {
	repos.PaginationQuery

	Id       uuid.UUID
	Codename string
	ClientId string
	State    string
	Type     string
	AnyId    string
}

type Application struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primary_key;"`
	CreatedAt          time.Time  `gorm:"type:timestamp"`
	UpdatedAt          time.Time  `gorm:"type:timestamp"`
	DeletedAt          *time.Time `gorm:"type:timestamp"`
	Codename           string     `gorm:"varchar"`
	Name               string     `gorm:"varchar"`
	Type               string     `gorm:"varchar"`
	State              string     `gorm:"varchar"`
	Description        string     `gorm:"varchar"`
	ClientId           string     `gorm:"varchar"`
	RedirectUrisStr    string     `gorm:"type:text;column:redirect_uris"`
	AvailableScopesStr string     `gorm:"type:text;column:available_scopes"`
}

func (a *Application) IsActive() bool {
	return a.State == "active"
}

func (a *Application) Scopes() []string {
	return strings.Split(a.AvailableScopesStr, ";")
}

func (a *Application) AddScopes(newScopes []string) {
	a.SetScopes(utils.StringSliceMerge(a.Scopes(), newScopes))
}

func (a *Application) SetScopes(scopes []string) {
	a.AvailableScopesStr = strings.Join(scopes, ";")
}

func (Application) TableName() string {
	return "applications"
}

func (a *Application) RedirectUris() []string {
	return strings.Split(a.RedirectUrisStr, "\n")
}

func (a *Application) AddUris(newUris []string) {
	a.SetUris(utils.StringSliceMerge(a.RedirectUris(), newUris))
}

func (a *Application) SetUris(uris []string) {
	a.RedirectUrisStr = strings.Join(uris, "\n")
}

type Repository interface {
	Create(ctx context.Context, app *Application) error
	Update(ctx context.Context, app *Application) error
	Delete(ctx context.Context, app *Application) error
	Query(ctx context.Context, query FindQuery) ([]Application, error)
	QueryOne(ctx context.Context, query FindQuery) (*Application, error)
}

func NweApplicationsRepositoryDB(db *gorm.DB) Repository {
	return &repositoryDB{
		DB:     db,
		common: repos.NewCommonRepositoryDB(db, "Applications"),
	}
}

type repositoryDB struct {
	DB     *gorm.DB
	common repos.CommonRepositoryDB
}

func (r *repositoryDB) Query(ctx context.Context, query FindQuery) (result []Application, err error) {
	db, entry := r.internalQueryBuilder(ctx, query)
	return result, r.common.ProcessQuery(db, &result, entry)
}

func (r *repositoryDB) QueryOne(ctx context.Context, query FindQuery) (*Application, error) {
	var result Application
	db, entry := r.internalQueryBuilder(ctx, query)
	one, err := r.common.ProcessQueryOne(db, &result, entry)
	if one == nil {
		return nil, err
	}
	return one.(*Application), err
}

func (r *repositoryDB) Create(ctx context.Context, app *Application) error {
	if app.ID == uuid.Nil {
		app.ID = uuid.NewV4()
	}
	return r.common.Create(ctx, app)
}

func (r *repositoryDB) Update(ctx context.Context, app *Application) error {
	return r.common.Update(ctx, app)
}

func (r *repositoryDB) Delete(ctx context.Context, app *Application) error {
	return r.common.Delete(ctx, app)
}

func (r *repositoryDB) internalQueryBuilder(ctx context.Context, query FindQuery) (*gorm.DB, *log.Entry) {
	db := r.DB
	logFields := log.Fields{
		"model": "application",
	}

	iid := uuid.FromStringOrNil(query.AnyId)

	if query.Id != uuid.Nil || iid != uuid.Nil {
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
	if query.AnyId != "" && iid == uuid.Nil {
		db = db.Where("codename = ? OR client_id = ?", query.AnyId, query.AnyId)
		logFields["anyid"] = query.AnyId
	}

	if query.State != "" {
		db = db.Where("state = ?", query.State)
	}

	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}

	db = r.common.AddPagination(db, logFields, query.PaginationQuery)

	return db, shared.GetLogger(ctx).WithFields(logFields)
}
