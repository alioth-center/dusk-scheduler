package init

import (
	"github.com/alioth-center/dusk-scheduler/app/config"
	"github.com/alioth-center/dusk-scheduler/infra/apis/location"
	db "github.com/alioth-center/dusk-scheduler/infra/database"
	"github.com/alioth-center/dusk-scheduler/infra/email"
	"github.com/alioth-center/dusk-scheduler/infra/sdk"
	"net/http"
)

func initDatabase(config *config.AppConfig) {
	databaseConfig := config.DatabaseConfig
	driverConfig := db.Config{
		StringSize:         uint(databaseConfig.StringSize),
		MaxIdleConnections: uint(databaseConfig.MaxIdleConnections),
		MaxOpenConnections: uint(databaseConfig.MaxOpenConnections),
		ConnectionLifeTime: uint(databaseConfig.ConnectionLifeTime),
		DriverOptions:      databaseConfig.DriverOptions,
		DataSource:         databaseConfig.DataSource,
	}

	var driver db.Database
	switch databaseConfig.Driver {
	case "mysql", "mariadb":
		driver = db.NewMySqlDatabase()
	case "postgres":
		driver = db.NewPostgresDatabase()
	case "sqlite":
		driver = db.NewSqliteDatabase()
	case "sqlserver":
		driver = db.NewSqlserverDatabase()
	default:
		panic("unsupported database driver")
	}

	instance, connectErr := db.ConnectDatabase(driver, databaseLogger, driverConfig)
	if connectErr != nil {
		panic(connectErr)
	}
	if !databaseConfig.MigrateDomains {
		database = instance

		return
	}

	migrateErr := instance.AutoMigrate(repositoryDomains...)
	if migrateErr != nil {
		panic(migrateErr)
	}

	database = instance
}

func initHttpClient(_ *config.AppConfig) {
	// todo: add http client options
	httpClient = http.DefaultClient
}

func initEmailSenderClient(config *config.AppConfig) {
	emailConfig := config.EmailConfig
	switch emailConfig.Provider {
	case "smtp":
		provider := emailConfig.SmtpProvider
		emailSenderClient = email.NewSmtpSenderClient(email.SmtpAuthSecret{
			Username: provider.Username,
			Password: provider.Password,
			Host:     provider.Host,
			Port:     uint16(provider.Port),
			Sender:   provider.Sender,
		})
	default:
		panic("unsupported email provider " + emailConfig.Provider)
	}
}

func initPositionLocator(config *config.AppConfig, httpClient *http.Client) {
	locatorConfig := config.PositionLocatorConfig
	switch locatorConfig.Provider {
	case "tencent_map":
		panic("tencent_map support not implemented yet")
	case "amap":
		panic("amap support not implemented yet")
	case "ip.sb":
		positionLocator = location.NewIpSbPositionLocator(httpClient)
	default:
		positionLocator = location.NewIpSbPositionLocator(httpClient)
	}
}

func initBrushSdk(_ *config.AppConfig, client *http.Client) {
	brushSdk = sdk.NewBrushSDK(client)
}
