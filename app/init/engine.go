package init

import (
	"github.com/alioth-center/dusk-scheduler/app/config"
	"github.com/alioth-center/dusk-scheduler/app/handler"
	"github.com/gin-gonic/gin"
)

func initGinEngine(config *config.AppConfig) {
	// todo: add gin engine options
	engineConfig := config.EngineConfig
	engine = gin.Default()

	handlers = []handler.Handler{brushHandler, clientHandler, outcomeHandler, painterHandler, taskHandler}
	if engineConfig.EnableAdminApi {

	}
}

func runGinEngine(config *config.AppConfig) {
	gin.SetMode(config.EngineConfig.RunMode)

	listenAt := config.EngineConfig.ListenAt
	if listenAt == "" {
		panic("invalid listening address")
	}

	_ = engine.Run(listenAt)
}
