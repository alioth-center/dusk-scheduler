package init

import (
	"github.com/alioth-center/dusk-scheduler/app/config"
	"github.com/alioth-center/dusk-scheduler/infra/email"
)

var (
	appConfig config.AppConfig
)

var (
	emailSenderClient email.SenderClient
)
