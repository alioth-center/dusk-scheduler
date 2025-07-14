package database

import (
	"fmt"
	"gorm.io/gorm"
)

type sqliteDatabase struct{}

func NewSqliteDatabase() Database { return &sqliteDatabase{} }

func (db *sqliteDatabase) DriverName() string {
	return "sqlite"
}

func (db *sqliteDatabase) DSN(dataSource map[string]any) string {
	return fmt.Sprintf("file:%s%s", dataSource["file"], dataSource["cache"])
}

func (db *sqliteDatabase) Config(_ map[string]any) *gorm.Config {
	// todo: support sqlite gorm config
	return &gorm.Config{}
}
