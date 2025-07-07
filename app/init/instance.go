package init

import (
	"github.com/alioth-center/dusk-scheduler/app/config"
	"github.com/alioth-center/dusk-scheduler/app/service"
	"github.com/alioth-center/dusk-scheduler/infra/apis/location"
	"github.com/alioth-center/dusk-scheduler/infra/email"
	"github.com/alioth-center/dusk-scheduler/infra/logger"
)

var (
	appConfig config.AppConfig
)

var (
	sysLogger         logger.Logger
	emailSenderClient email.SenderClient
	positionLocator   location.PositionLocator
)

var (
	emailService    service.EmailService
	locationService service.LocationService
)
