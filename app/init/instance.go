package init

import (
	"github.com/alioth-center/dusk-scheduler/app/config"
	"github.com/alioth-center/dusk-scheduler/app/repository"
	"github.com/alioth-center/dusk-scheduler/app/service"
	"github.com/alioth-center/dusk-scheduler/infra/apis/location"
	"github.com/alioth-center/dusk-scheduler/infra/cache"
	"github.com/alioth-center/dusk-scheduler/infra/email"
	"github.com/alioth-center/dusk-scheduler/infra/logger"
	"gorm.io/gorm"
	"net/http"
)

var (
	httpClient *http.Client
	appConfig  config.AppConfig
	database   *gorm.DB
	caching    cache.Cache
)

var (
	sysLogger         logger.Logger
	emailSenderClient email.SenderClient
	positionLocator   location.PositionLocator
)

var (
	clientDao      repository.ClientDao
	taskDao        repository.TaskDao
	promotionalDao repository.PromotionalDao
	painterDao     repository.PainterDao
	outcomeDao     repository.OutcomeDao
	storageDao     repository.StorageDao

	authorizationCache    repository.AuthorizationCache
	quotaCache            repository.QuotaCache
	painterHeartbeatCache repository.PainterHeartbeatCache
)

var (
	emailService    service.EmailService
	locationService service.LocationService
	clientService   service.ClientService
	outcomeService  service.OutcomeService
	painterService  service.PainterService
)
