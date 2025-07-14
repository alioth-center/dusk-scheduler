package database

import (
	"fmt"
	"gorm.io/gorm"
)

type postgresDatabase struct{}

func NewPostgresDatabase() Database { return &postgresDatabase{} }

func (db *postgresDatabase) DriverName() string {
	return "postgres"
}

func (db *postgresDatabase) DSN(dataSource map[string]any) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		dataSource["host"], dataSource["port"],
		dataSource["username"], dataSource["password"],
		dataSource["dbname"], dataSource["ssl_mode"], dataSource["location"],
	)
}

func (db *postgresDatabase) Config(_ map[string]any) *gorm.Config {
	// todo: support postgres gorm config
	return &gorm.Config{}
}
