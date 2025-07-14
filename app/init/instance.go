package init

import (
	"github.com/alioth-center/dusk-scheduler/app/config"
	"github.com/alioth-center/dusk-scheduler/app/domain"
	"github.com/alioth-center/dusk-scheduler/app/handler"
	"github.com/alioth-center/dusk-scheduler/app/repository"
	"github.com/alioth-center/dusk-scheduler/app/service"
	"github.com/alioth-center/dusk-scheduler/infra/apis/location"
	"github.com/alioth-center/dusk-scheduler/infra/cache"
	"github.com/alioth-center/dusk-scheduler/infra/email"
	"github.com/alioth-center/dusk-scheduler/infra/logger"
	"github.com/alioth-center/dusk-scheduler/infra/sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

var (
	appConfig config.AppConfig
)

var (
	httpClient        *http.Client
	sysLogger         logger.Logger
	databaseLogger    logger.DatabaseLogger
	emailSenderClient email.SenderClient
	positionLocator   location.PositionLocator
	brushSdk          sdk.BrushSDK
	database          *gorm.DB
	caching           cache.Cache
	engine            *gin.Engine
)

var (
	clientDao      repository.ClientDao
	taskDao        repository.TaskDao
	promotionalDao repository.PromotionalDao
	painterDao     repository.PainterDao
	outcomeDao     repository.OutcomeDao
	storageDao     repository.StorageDao
	brushDao       repository.BrushDao

	authorizationCache    repository.AuthorizationCache
	quotaCache            repository.QuotaCache
	painterHeartbeatCache repository.PainterHeartbeatCache
	taskContentCache      repository.TaskContentCache
	brushCache            repository.BrushCache
)

var (
	emailService       service.EmailService
	locationService    service.LocationService
	clientService      service.ClientService
	outcomeService     service.OutcomeService
	painterService     service.PainterService
	promotionalService service.PromotionalService
	brushService       service.BrushService
	taskService        service.TaskService
)

var (
	brushHandler   *handler.BrushHandler
	clientHandler  *handler.ClientHandler
	outcomeHandler *handler.OutcomeHandler
	painterHandler *handler.PainterHandler
	taskHandler    *handler.TaskHandler
)

var (
	handlers []handler.Handler

	repositoryDomains = []any{
		&domain.Brush{}, &domain.Client{}, &domain.Outcome{}, &domain.Painter{}, &domain.Promotional{},
		&domain.Storage{}, &domain.Task{},
	}
)
