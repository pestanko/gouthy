package users

import (
	"context"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/shared/repos"
	"github.com/pestanko/gouthy/app/shared/utils"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type FindQuery struct {
	repos.PaginationQuery
	Id       uuid.UUID
	Username string
	Email    string
	AnyId    string
}

type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time  `gorm:"type:timestamp"`
	UpdatedAt time.Time  `gorm:"type:timestamp"`
	DeletedAt *time.Time `gorm:"type:timestamp"`
	Username  string
	Password  string
	Name      string
	Email     string
	State     string
	UserType  string     `gorm:"column:user_type"`
}

func (User) TableName() string {
	return "users"
}

func (user *User) SetPassword(password string) error {
	hash, err := utils.HashString(password)
	if err != nil {
		return err
	}

	user.Password = hash
	return nil
}

func (user *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
}

func (user *User) ToEntity() *UserDTO {
	return &UserDTO{
		baseUserDTO: *convertModelToUserBase(user),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func (user *User) LogFields() log.Fields {
	return log.Fields{
		"user_id":  user.ID,
		"username": user.Username,
	}
}

type Repository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, user *User) error
	Query(ctx context.Context, query FindQuery) ([]User, error)
	QueryOne(ctx context.Context, query FindQuery) (*User, error)
}

type repositoryDB struct {
	DB     *gorm.DB
	common repos.CommonRepositoryDB
}

func (r *repositoryDB) Create(ctx context.Context, user *User) error {
	return r.common.Create(ctx, user)
}

func (r *repositoryDB) Update(ctx context.Context, user *User) error {
	return r.common.Update(ctx, user)
}

func (r *repositoryDB) Delete(ctx context.Context, user *User) error {
	return r.common.Delete(ctx, user)
}

func (r *repositoryDB) QueryOne(ctx context.Context, query FindQuery) (*User, error) {
	var result User
	db, entry := r.internalQueryBuilder(ctx, query)
	one, err := r.common.ProcessQueryOne(db, &result, entry)
	if one == nil {
		return nil, err
	}
	return one.(*User), err

}

func (r *repositoryDB) Query(ctx context.Context, query FindQuery) (result []User, err error) {
	db, entry := r.internalQueryBuilder(ctx, query)
	return result, r.common.ProcessQuery(db, &result, entry)
}

func (r *repositoryDB) internalQueryBuilder(ctx context.Context, query FindQuery) (*gorm.DB, *log.Entry) {
	db := r.DB
	logFields := log.Fields{
		"model": "user",
	}

	iid := uuid.FromStringOrNil(query.AnyId)

	if iid != uuid.Nil {
		query.Id = iid
	}

	if query.Id != uuid.Nil {
		db = db.Where("id = ?", query.Id)
		logFields["id"] = query.Id
	}
	if query.Email != "" {
		db = db.Where("email = ?", query.Email)
		logFields["email"] = query.Email
	}

	if query.Username != "" {
		db = db.Where("username = ?", query.Username)
		logFields["username"] = query.Username
	}

	if query.AnyId != "" && iid == uuid.Nil {
		db = db.Where("username = ?", query.AnyId)
		logFields["username"] = query.AnyId
	}

	db = r.common.AddPagination(db, logFields, query.PaginationQuery)

	return db, shared.GetLogger(ctx).WithFields(logFields)
}

func NewUsersRepositoryDB(db *gorm.DB) Repository {
	return &repositoryDB{DB: db, common: repos.NewCommonRepositoryDB(db, "User")}
}
