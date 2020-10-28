package shared

import (
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
)

// GetDBConnection - Get the database connection
func GetDBConnection(config *AppConfig) (DBConnection, error) {
	dbConfig := config.DB.GetDefault()
	dialector := getOpenerBasedOnConfig(dbConfig)
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	return newGormConnection(db), nil
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

func getOpenerBasedOnConfig(dbConfig *DBEntryConfig) gorm.Dialector {
	switch strings.ToLower(dbConfig.DBType) {
	case "sqlite", "memory":
		return sqlite.Open(dbConfig.Uri)
	case "postgres":
		return postgres.Open(dbConfig.Uri)
	}
	return sqlite.Open(dbConfig.Uri)
}
