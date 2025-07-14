package init

import (
	"github.com/alioth-center/dusk-scheduler/app/handler"
	"github.com/alioth-center/dusk-scheduler/app/repository"
	"github.com/alioth-center/dusk-scheduler/app/service"
)

func init() {
	initConfig()
	initInfra()
	initRepository()
	initService()
	initHandler()
	initEngine()
}

func initInfra() {
	initHttpClient(&appConfig)
	initEmailSenderClient(&appConfig)
	initPositionLocator(&appConfig, httpClient)
	initBrushSdk(&appConfig, httpClient)
	initDatabase(&appConfig)
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
	taskService = service.NewTaskService(taskDao, outcomeDao, taskContentCache, sysLogger, &appConfig)
	painterService = service.NewPainterService(storageDao, painterDao, painterHeartbeatCache, sysLogger, locationService, &appConfig)
	brushService = service.NewBrushService(brushDao, taskDao, brushCache, taskContentCache, brushSdk, sysLogger)
}

func initHandler() {
	brushHandler = handler.NewBrushHandler(emailService, brushService)
	clientHandler = handler.NewClientHandler(emailService, clientService, taskService, promotionalService)
	outcomeHandler = handler.NewOutcomeHandler(taskService, outcomeService, painterService)
	painterHandler = handler.NewPainterHandler(taskService, outcomeService, emailService, painterService)
	taskHandler = handler.NewTaskHandler(clientService, taskService, outcomeService, brushService)
}

func initEngine() {
	initGinEngine(&appConfig)

	routerGroup := engine.Group("/dusk-scheduler")
	for _, h := range handlers {
		h.RegisterHandler(routerGroup)
	}
}

func RunApp() {
	runGinEngine(&appConfig)
}
