package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/alioth-center/dusk-scheduler/app/domain"
	"github.com/alioth-center/dusk-scheduler/app/repository"
	"github.com/alioth-center/dusk-scheduler/app/service/errors"
	"github.com/alioth-center/dusk-scheduler/infra/logger"
	"github.com/alioth-center/dusk-scheduler/infra/sdk"
	"time"
)

type brushService struct {
	brushDao         repository.BrushDao
	taskDao          repository.TaskDao
	brushCache       repository.BrushCache
	taskContentCache repository.TaskContentCache
	brushSdk         sdk.BrushSDK
	sysLogger        logger.Logger
}

func NewBrushService(
	brushDao repository.BrushDao,
	taskDao repository.TaskDao,
	brushCache repository.BrushCache,
	taskContentCache repository.TaskContentCache,
	brushSdk sdk.BrushSDK,
	sysLogger logger.Logger,
) BrushService {
	return &brushService{
		brushDao:         brushDao,
		taskDao:          taskDao,
		brushCache:       brushCache,
		taskContentCache: taskContentCache,
		brushSdk:         brushSdk,
		sysLogger:        sysLogger,
	}
}

func (srv *brushService) CreateBrush(ctx context.Context, maintainer string, protocol string, callURL string) (brushID uint64, err error) {
	var protocolEnum domain.BrushProtocol
	switch protocol {
	case domain.BrushProtocolHttp.String():
		protocolEnum = domain.BrushProtocolHttp
	case domain.BrushProtocolGrpc.String():
		protocolEnum = domain.BrushProtocolGrpc
	case domain.BrushProtocolTcp.String():
		protocolEnum = domain.BrushProtocolTcp
	default:
		return 0, errors.RegisterBrushProtocolNotSupportError()
	}

	brush := domain.Brush{
		Protocol:    protocolEnum,
		Maintainer:  maintainer,
		CallURL:     callURL,
		RegisterAt:  time.Now(),
		ConnectedAt: time.Now(),
	}
	id, createErr := srv.brushDao.CreateBrush(ctx, &brush)
	if createErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to create brush: [%s] %s", protocol, maintainer), createErr)

		return 0, createErr
	}
	if cacheErr := srv.brushCache.AddBrush(ctx, id); cacheErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to cache brush: %d", id), cacheErr)

		return 0, cacheErr
	}

	return id, nil
}

func (srv *brushService) DisconnectBrush(ctx context.Context, brushID uint64) (err error) {
	if deleteErr := srv.brushDao.UpdateBrushAsDisconnected(ctx, brushID); deleteErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to disconnect brush: %d", brushID), deleteErr)

		return deleteErr
	}
	if cacheErr := srv.brushCache.RemoveBrush(ctx, brushID); cacheErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to remove brush from cache: %d", brushID), cacheErr)

		return cacheErr
	}

	return nil
}

func (srv *brushService) RenderImage(ctx context.Context, taskID uint64) (result *bytes.Buffer, err error) {
	defer func() {
		reason := domain.TaskArchiveReasonAcknowledged
		if result != nil {
			reason = domain.TaskArchiveReasonFailed
		}

		archiveErr := srv.taskDao.UpdateTaskAsArchived(ctx, taskID, reason)
		if archiveErr != nil {
			srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to archive task: %d", taskID), archiveErr)

			result, err = nil, archiveErr
		}
	}()

	picked, pickErr := srv.brushCache.GetRandomBrush(ctx)
	if pickErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to pick random brush for task: %d", taskID), pickErr)

		return nil, pickErr
	}
	if picked == 0 {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("no available brush for task: %d", taskID), nil)

		return nil, errors.RenderImageNoAvailableBrushError()
	}
	instance, existInstance, getInstanceErr := srv.brushDao.GetBrushByID(ctx, picked)
	if getInstanceErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get brush by id: %d", picked), getInstanceErr)

		return nil, getInstanceErr
	}
	if !existInstance {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("brush not found: %d", picked), nil)

		return nil, errors.RenderImageBrushNotFoundError()
	}

	encoded, getContentErr := srv.taskContentCache.GetTaskContent(ctx, taskID)
	if getContentErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get task content for task: %d", taskID), getContentErr)

		return nil, getContentErr
	}
	result, err = srv.brushSdk.RenderImage(ctx, instance.Protocol.String(), instance.CallURL, encoded)
	if err != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to render image for task: %d with brush: %d", taskID, picked), err)

		return nil, err
	}

	// todo: update quota usage for brush

	return result, err
}
