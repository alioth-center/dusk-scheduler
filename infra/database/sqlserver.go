package database

import (
	"fmt"
	"gorm.io/gorm"
)

type sqlserverDatabase struct{}

func NewSqlserverDatabase() Database { return &sqlserverDatabase{} }

func (db *sqlserverDatabase) DriverName() string {
	return "sqlserver"
}

func (db *sqlserverDatabase) DSN(dataSource map[string]any) string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		dataSource["username"], dataSource["password"],
		dataSource["host"], dataSource["port"], dataSource["database"],
	)
}

func (db *sqlserverDatabase) Config(_ map[string]any) *gorm.Config {
	// todo: support sqlserver gorm config
	return &gorm.Config{}
}
