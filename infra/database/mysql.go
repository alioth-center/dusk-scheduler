package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlDatabase struct{}

func NewMySqlDatabase() Database { return &mysqlDatabase{} }

func (db *mysqlDatabase) DriverName() string {
	return mysql.DefaultDriverName
}

func (db *mysqlDatabase) DSN(dataSource map[string]any) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s",
		dataSource["username"], dataSource["password"],
		dataSource["host"], dataSource["port"],
		dataSource["database"], dataSource["charset"], dataSource["location"],
	)
}

func (db *mysqlDatabase) Config(_ map[string]any) *gorm.Config {
	// todo: support mysql gorm config
	return &gorm.Config{}
}
