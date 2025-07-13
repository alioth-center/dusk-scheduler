package init

import (
	"github.com/alioth-center/dusk-scheduler/app/config"
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
	database  *gorm.DB
	caching   cache.Cache
	engine    *gin.Engine
)

var (
	httpClient        *http.Client
	sysLogger         logger.Logger
	emailSenderClient email.SenderClient
	positionLocator   location.PositionLocator
	brushSdk          sdk.BrushSDK
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

	handlers []handler.Handler
)
