package database

import (
	"errors"
	"gorm.io/driver/mysql"
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
		driver, err = gorm.Open(mysql.Open(db.DSN()), &gorm.Config{})
	default:

	}

	return driver, err
}
