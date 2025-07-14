package database

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type Database interface {
	DriverName() string
	DSN(dataSource map[string]any) string
	Config(options map[string]any) *gorm.Config
}

type Config struct {
	StringSize         uint
	MaxIdleConnections uint
	MaxOpenConnections uint
	ConnectionLifeTime uint
	DriverName         string
	DriverOptions      map[string]any
	DataSource         map[string]any
}

func ConnectDatabase(db Database, logger logger.Interface, conf Config) (driver *gorm.DB, err error) {
	if db == nil {
		err = errors.New("database is nil")

		return nil, err
	}

	switch conf.DriverName {
	case "mysql":
		driver, err = gorm.Open(mysql.New(mysql.Config{
			DriverName:        db.DriverName(),
			DSN:               db.DSN(conf.DataSource),
			DefaultStringSize: conf.StringSize,
		}), db.Config(conf.DriverOptions))
	case "mariadb":
		driver, err = gorm.Open(mysql.New(mysql.Config{
			DriverName:                db.DriverName(),
			DSN:                       db.DSN(conf.DataSource),
			SkipInitializeWithVersion: false,
			DefaultStringSize:         conf.StringSize,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
		}), db.Config(conf.DriverOptions))
	case "postgres":
		driver, err = gorm.Open(postgres.New(postgres.Config{
			DriverName: db.DriverName(),
			DSN:        db.DSN(conf.DataSource),
		}), db.Config(conf.DriverOptions))
	case "sqlite":
		driver, err = gorm.Open(sqlite.Open(db.DSN(conf.DataSource)), db.Config(conf.DriverOptions))
	case "sqlserver":
		driver, err = gorm.Open(sqlserver.Open(db.DSN(conf.DataSource)), db.Config(conf.DriverOptions))
	default:
		panic("unsupported database driver")
	}
	if err != nil {
		return nil, err
	}

	driver.Logger = logger
	rawDB, getErr := driver.DB()
	if getErr != nil {
		return nil, getErr
	}
	rawDB.SetMaxIdleConns(int(conf.MaxIdleConnections))
	rawDB.SetMaxOpenConns(int(conf.MaxOpenConnections))
	rawDB.SetConnMaxLifetime(time.Duration(conf.ConnectionLifeTime) * time.Second)

	return driver, nil
}
