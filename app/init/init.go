package init

import (
	"github.com/alioth-center/dusk-scheduler/app/repository"
	"github.com/alioth-center/dusk-scheduler/app/service"
	"github.com/alioth-center/dusk-scheduler/infra/config"
	"os"
)

func init() {
	initConfig()
	initInfra()
	initRepository()
	initService()
	initHandler()
}

func initConfig() {
	switch os.Getenv("CONFIG_TYPE") {
	case "file":
		configReader := config.NewFileConfig(os.Getenv("CONFIG_FILE_PATH"))
		if readErr := configReader.ParseAppConfig("", "", &appConfig); readErr != nil {
			panic(readErr)
		}
	case "apollo":
		panic("apollo support not implemented yet")
	case "remote_url":
		configReader := config.NewRemoteURLConfig()
		if readErr := configReader.ParseAppConfig(os.Getenv("CONFIG_URL_PATH"), "", &appConfig); readErr != nil {
			panic(readErr)
		}
	default:
		panic("invalid config type")
	}
}

func initInfra() {
	initEmailSenderClient(&appConfig)
	initPositionLocator(&appConfig, httpClient)
}

func initRepository() {
	clientDao = repository.NewClientDao(database)
	taskDao = repository.NewTaskDao(database)
	promotionalDao = repository.NewPromotionalDao(database)

	authorizationCache = repository.NewAuthorizationCache(caching)
	quotaCache = repository.NewQuotaCache(caching)
}

func initService() {
	emailService = service.NewEmailService(emailSenderClient, sysLogger, &appConfig)
	locationService = service.NewLocationService(positionLocator, sysLogger)
	clientService = service.NewClientService(clientDao, promotionalDao, taskDao, authorizationCache, quotaCache, locationService, sysLogger, &appConfig)
	outcomeService = service.NewOutcomeService(taskDao, painterDao, outcomeDao, storageDao, sysLogger, httpClient)
	painterService = service.NewPainterService(storageDao, painterDao, painterHeartbeatCache, sysLogger, locationService, &appConfig)
}

func initHandler() {

}
