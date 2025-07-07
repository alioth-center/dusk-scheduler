package service

import (
	"context"
	"fmt"
	"github.com/alioth-center/dusk-scheduler/infra/apis/location"
	"github.com/alioth-center/dusk-scheduler/infra/logger"
)

type locationService struct {
	positionLocator location.PositionLocator
	sysLogger       logger.Logger
}

func NewLocationService(
	positionLocator location.PositionLocator,
	sysLogger logger.Logger,
) LocationService {
	return &locationService{
		positionLocator: positionLocator,
		sysLogger:       sysLogger,
	}
}

func (srv *locationService) DetectIPLocation(ctx context.Context, ip string) (address *location.Address, err error) {
	detected, detectErr := srv.positionLocator.DetectIP(ctx, ip)
	if detectErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to detect ip location: %s", ip), detectErr)

		return nil, detectErr
	}

	return detected, nil
}
