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

type painterService struct {
	storageDao            repository.StorageDao
	painterDao            repository.PainterDao
	painterHeartbeatCache repository.PainterHeartbeatCache
	sysLogger             logger.Logger
	locationService       LocationService
	appConfig             *config.AppConfig
}

func NewPainterService(
	storageDao repository.StorageDao,
	painterDao repository.PainterDao,
	painterHeartbeatCache repository.PainterHeartbeatCache,
	sysLogger logger.Logger,
	locationService LocationService,
	appConfig *config.AppConfig,
) PainterService {
	return &painterService{
		storageDao:            storageDao,
		painterDao:            painterDao,
		painterHeartbeatCache: painterHeartbeatCache,
		sysLogger:             sysLogger,
		locationService:       locationService,
		appConfig:             appConfig,
	}
}

func (srv *painterService) CreatePainter(ctx context.Context, maintainer string, slot int, ip string) (painter *domain.Painter, policy *domain.Storage, err error) {
	options := srv.appConfig.PainterOptions
	policyEntity, existPolicy, getPolicyErr := srv.storageDao.GetStorageByName(ctx, options.StoragePolicy)
	if getPolicyErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get policy by name: %s", options.StoragePolicy), getPolicyErr)

		return nil, nil, getPolicyErr
	}
	if !existPolicy {
		srv.sysLogger.WarnCtx(ctx, fmt.Sprintf("policy does not exist: %s", options.StoragePolicy), nil)

		return nil, nil, errors.RegisterPainterStoragePolicyNotFoundError()
	}

	location, detectErr := srv.locationService.DetectIPLocation(ctx, ip)
	if detectErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to detect ip location: %s", ip), detectErr)

		return nil, nil, detectErr
	}
	position := location.Region
	if len(location.City) > 0 {
		position = fmt.Sprintf("%s %s", location.Region, location.City)
	}

	painterEntity := domain.Painter{
		Location:    position,
		Maintainer:  maintainer,
		Secret:      utils.GenerateToken(32, ""),
		Slot:        uint8(slot),
		RegisterAt:  time.Now(),
		ConnectedAt: time.Now(),
		ConnectTime: 1,
		PolicyID:    policyEntity.ID,
	}
	painterID, createErr := srv.painterDao.CreatePainter(ctx, &painterEntity)
	if createErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to create painter: %s-%s", maintainer, ip), createErr)

		return nil, nil, createErr
	}

	var painterName string
	switch options.NamingRule {
	case "default":
		painterName = utils.GenerateNameByDefaultDictionary(painterID)
	default:
		if len(options.NamingDictionary[options.NamingRule]) == 0 {
			return nil, nil, errors.RegisterPainterInvalidNamingRuleError()
		}

		painterName = utils.GenerateName(painterID, options.NamingDictionary[options.NamingRule])
	}

	if updateErr := srv.painterDao.UpdatePainterName(ctx, painterID, painterName); updateErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to update painter name: %d", painterID), updateErr)

		return nil, nil, updateErr
	}
	if heartbeatErr := srv.painterHeartbeatCache.UpdateHeartbeatTime(ctx, painterName); heartbeatErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to update heartbeat time: %d", painterID), heartbeatErr)

		return nil, nil, heartbeatErr
	}
	painterEntity.ID, painterEntity.Name = painterID, painterName

	return &painterEntity, policy, nil
}

func (srv *painterService) ReconnectPainter(ctx context.Context, name string) (heartbeat bool, err error) {
	painter, existPainter, getPainterErr := srv.painterDao.GetPainterByName(ctx, name)
	if getPainterErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get painter by name: %s", name), getPainterErr)

		return false, getPainterErr
	}
	if !existPainter {
		srv.sysLogger.WarnCtx(ctx, fmt.Sprintf("painter does not exist: %s", name), nil)

		return false, errors.ReconnectPainterNotFoundError()
	}

	if heartbeatErr := srv.painterHeartbeatCache.UpdateHeartbeatTime(ctx, name); heartbeatErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to update heartbeat: %d", painter.ID), heartbeatErr)

		return false, heartbeatErr
	}

	if painter.ConnectedAt.After(painter.DisconnectedAt) {
		return true, nil
	}

	if updateErr := srv.painterDao.UpdatePainterAsConnected(ctx, painter.ID); updateErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to update painter as connected: %s", name), updateErr)

		return false, updateErr
	}

	return false, nil
}

func (srv *painterService) DisconnectPainter(ctx context.Context, name string) (err error) {
	painter, existPainter, getPainterErr := srv.painterDao.GetPainterByName(ctx, name)
	if getPainterErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get painter by name: %s", name), getPainterErr)

		return getPainterErr
	}
	if !existPainter {
		srv.sysLogger.WarnCtx(ctx, fmt.Sprintf("painter does not exist: %s", name), nil)

		return errors.DisconnectPainterNotFoundError()
	}

	if updateErr := srv.painterDao.UpdatePainterAsDisconnected(ctx, painter.ID); updateErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to update painter as disconnected: %s", name), updateErr)

		return updateErr
	}
	if heartbeatErr := srv.painterHeartbeatCache.DeleteHeartbeatTime(ctx, name); heartbeatErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to delete heartbeat: %d", painter.ID), heartbeatErr)

		return heartbeatErr
	}

	return nil
}

func (srv *painterService) GetPainterByID(ctx context.Context, painterID uint64) (painter *domain.Painter, exist bool, err error) {
	painter, exist, err = srv.painterDao.GetPainterByID(ctx, painterID)
	if err != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get painter by id: %d", painterID), err)

		return nil, false, err
	}

	return painter, exist, nil
}
