package shared

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewInMemoryConnection() (DBConnection, error) {
	cxn := "file:memdb1?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(cxn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &gormConnection{
		db: db,
	}, nil
}

func NewPostgresConnection(config *AppConfig) (DBConnection, error) {
	cxn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
		config.DB.Host, config.DB.Port, config.DB.User, config.DB.DBName, config.DB.Password, config.DB.SSLMode)
	db, err := gorm.Open(postgres.Open(cxn), &gorm.Config{})
	return newGormConnection(db), err
}

// GetDBConnection - Get the database connection
func GetDBConnection(config *AppConfig) (DBConnection, error) {
	if config.DB.InMemory {
		log.WithField("db_type", "in-memory").Info("Starting with database")
		return NewInMemoryConnection()
	} else {
		log.WithField("db_type", "postgres").Info("Starting with database")
		return NewPostgresConnection(config)
	}
}

func DBConnectionIntoGorm(db DBConnection) *gorm.DB {
	raw := db.GetRaw()
	switch raw.(type) {
	case *gorm.DB:
		return raw.(*gorm.DB)
	default:
		return nil
	}
}

type DBConnection interface {
	GetRaw() interface{}
	Close() error
}

type gormConnection struct {
	db *gorm.DB
}

func (g *gormConnection) GetRaw() interface{} {
	return g.db
}

func (g *gormConnection) Close() error {
	return nil
}

func newGormConnection(db *gorm.DB) DBConnection {
	return &gormConnection{
		db,
	}
}
