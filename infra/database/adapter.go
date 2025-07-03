package database

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type Database interface {
	DriverName() string
	DSN() string
	Config() *gorm.Config
}

func ConnectDatabase(db Database) (driver *gorm.DB, err error) {
	if db == nil {
		err = errors.New("database is nil")

		return nil, err
	}

	switch db.DriverName() {
	case "mysql":
		driver, err = gorm.Open(mysql.Open(db.DSN()), db.Config())
	case "postgres":
		driver, err = gorm.Open(postgres.Open(db.DSN()), db.Config())
	case "sqlite":
		driver, err = gorm.Open(sqlite.Open(db.DSN()), db.Config())
	case "sqlserver":
		driver, err = gorm.Open(sqlserver.Open(db.DSN()), db.Config())
	default:
		panic("unsupported database driver")
	}

	return driver, err
}
