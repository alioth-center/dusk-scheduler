package init

import (
	"github.com/alioth-center/dusk-scheduler/app/config"
	"github.com/gin-gonic/gin"
)

func initGinEngine(_ *config.AppConfig) {
	// todo: add gin engine options
	engine = gin.Default()
}

func runGinEngine(config *config.AppConfig) {
	gin.SetMode(config.EngineConfig.RunMode)

	listenAt := config.EngineConfig.ListenAt
	if listenAt == "" {
		panic("invalid listening address")
	}

	engine.Run(listenAt)
}
