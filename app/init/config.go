package init

import (
	"github.com/alioth-center/dusk-scheduler/infra/config"
	"os"
)

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
