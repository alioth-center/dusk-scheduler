package service

import (
	"context"
	"fmt"
	"github.com/alioth-center/dusk-scheduler/app/config"
	"github.com/alioth-center/dusk-scheduler/app/domain"
	"github.com/alioth-center/dusk-scheduler/app/repository"
	"github.com/alioth-center/dusk-scheduler/app/service/errors"
	"github.com/alioth-center/dusk-scheduler/infra/logger"
	"github.com/alioth-center/dusk-scheduler/infra/utils"
	"time"
)

type clientService struct {
	clientDao          repository.ClientDao
	promotionalDao     repository.PromotionalDao
	taskDao            repository.TaskDao
	authorizationCache repository.AuthorizationCache
	quotaCache         repository.QuotaCache
	locationService    LocationService
	sysLogger          logger.Logger
	appConfig          *config.AppConfig
}

func NewClientService(
	clientDao repository.ClientDao,
	promotionalDao repository.PromotionalDao,
	taskDao repository.TaskDao,
	authorizationCache repository.AuthorizationCache,
	quotaCache repository.QuotaCache,
	locationService LocationService,
	sysLogger logger.Logger,
	appConfig *config.AppConfig,
) ClientService {
	return &clientService{
		clientDao:          clientDao,
		promotionalDao:     promotionalDao,
		taskDao:            taskDao,
		authorizationCache: authorizationCache,
		quotaCache:         quotaCache,
		locationService:    locationService,
		sysLogger:          sysLogger,
		appConfig:          appConfig,
	}
}

func (srv *clientService) CreateClient(ctx context.Context, email string, promotionCode string, ip string) (client *domain.Client, err error) {
	if len(promotionCode) == 0 {
		promotionCode = srv.appConfig.ClientOptions.DefaultPermission.PromotionalCode
	}

	// find client permission
	promotional, existPromotional, queryPromotionalErr := srv.promotionalDao.GetPromotionalByCode(ctx, promotionCode)
	if queryPromotionalErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get promotional by code: %s", promotionCode), queryPromotionalErr)

		return nil, queryPromotionalErr
	} else if !existPromotional {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("existance promotion code not found: %s", promotionCode), nil)

		return nil, errors.RedemptionCodeNotFoundError()
	}

	// detect client ip position
	location, detectErr := srv.locationService.DetectIPLocation(ctx, ip)
	if detectErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get location: %s", ip), detectErr)

		return nil, detectErr
	}
	position := location.Region
	if len(location.City) > 0 {
		position = fmt.Sprintf("%s,%s", location.Region, location.City)
	}

	// initialize client entity and insert into database
	clientEntity := domain.Client{
		Maintainer:      email,
		Region:          position,
		ApiKey:          utils.GenerateToken(32, srv.appConfig.ClientOptions.ClientApiKeyPrefix),
		AllowBrush:      promotional.AllowBrush,
		AllowDelay:      promotional.AllowDelay,
		AllowHeight:     promotional.AllowHeight,
		AllowWidth:      promotional.AllowWidth,
		AllowPriority:   promotional.AllowPriority,
		Quota:           promotional.Quota,
		LimitRenderSize: promotional.LimitRenderSize,
		LimitFrequency:  promotional.LimitFrequency,
		LimitDuration:   promotional.LimitDuration,
	}
	clientID, createErr := srv.clientDao.CreateClient(ctx, &clientEntity)
	if createErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to create client for email: %s", email), createErr)

		return nil, createErr
	}
	clientEntity.ID = clientID

	return &clientEntity, nil
}

func (srv *clientService) StoreAuthorizationCode(ctx context.Context, clientID uint64, authorizationCode string) (expiredAt time.Time, err error) {
	expire := time.Duration(srv.appConfig.ClientOptions.AuthCodeExpireSeconds) * time.Second
	if cacheErr := srv.authorizationCache.StoreAuthorizationCode(ctx, clientID, authorizationCode, expire); cacheErr != nil {
		return time.Time{}, fmt.Errorf("failed to store authorization code for client %d: %w", clientID, cacheErr)
	}

	return time.Now().Add(expire), nil
}

func (srv *clientService) AuthorizeClient(ctx context.Context, clientID uint64, emailAddress string, authorizationCode string) (authorized bool, maintainer string, apiKey string, err error) {
	if len(authorizationCode) != 6 {
		return false, "", "", errors.InvalidAuthorizationCodeError()
	}

	// check authorization code
	code, exist, cacheErr := srv.authorizationCache.GetAuthorizationCode(ctx, clientID)
	if cacheErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get authorization code for client %d", clientID), cacheErr)

		return false, "", "", cacheErr
	}
	if !exist || code != authorizationCode {
		return false, "", "", errors.EmailAddressOrCodeMismatchError()
	}

	// check maintainer email address and client id
	client, existClient, queryErr := srv.clientDao.GetClientByID(ctx, clientID)
	if queryErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get client by id: %d", clientID), queryErr)

		return false, "", "", queryErr
	}
	if !existClient {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("existance client not found by id: %d", clientID), nil)

		return false, "", "", errors.EmailAddressOrCodeMismatchError()
	}
	if client.Maintainer != emailAddress {
		return false, "", "", errors.EmailAddressOrCodeMismatchError()
	}

	// desensitize email address
	conf := srv.appConfig.ClientOptions
	maintainerEmail := utils.DesensitizeEmailAddress(client.Maintainer, conf.DesensitizePrefix, conf.DesensitizeSuffix, conf.DesensitizeDomain)

	return true, maintainerEmail, client.ApiKey, nil
}

func (srv *clientService) GetClientData(ctx context.Context, clientID uint64) (client *domain.Client, exist bool, err error) {
	client, exist, err = srv.clientDao.GetClientByID(ctx, clientID)
	if err != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get client by id: %d", clientID), err)

		return nil, false, err
	}

	return client, exist, nil
}

func (srv *clientService) GetClientQuotaUsage(ctx context.Context, clientID uint64) (quotaTotal, quotaUsage uint64, lastCheckTime time.Time, err error) {
	checkAt, getErr := srv.quotaCache.LastStatisticsAt(ctx)
	if getErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get last statistics time for client %d", clientID), getErr)

		return 0, 0, time.Time{}, getErr
	}

	totalQuota, quotaErr := srv.quotaCache.GetTotalQuota(ctx, clientID)
	if quotaErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get total quota for client %d", clientID), quotaErr)

		return 0, 0, time.Time{}, quotaErr
	}

	currentUsage, usageErr := srv.taskDao.StatisticsClientQuotaUsage(ctx, clientID, checkAt)
	if usageErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get quota usage for client %d", clientID), usageErr)

		return 0, 0, time.Time{}, usageErr
	}

	return totalQuota, currentUsage, checkAt, nil
}
