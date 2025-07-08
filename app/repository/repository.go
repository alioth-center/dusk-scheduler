package repository

import (
	"context"
	"github.com/alioth-center/dusk-scheduler/app/domain"
	"time"
)

type ClientDao interface {
	CreateClient(ctx context.Context, client *domain.Client) (clientID uint64, err error)
	GetClientByID(ctx context.Context, clientID uint64) (client *domain.Client, exist bool, err error)
}

type TaskDao interface {
	StatisticsClientQuotaUsage(ctx context.Context, clientID uint64, startTime time.Time) (usage uint64, err error)
}

type PromotionalDao interface {
	GetPromotionalByCode(ctx context.Context, code string) (promotional *domain.Promotional, exist bool, err error)
}

type AuthorizationCache interface {
	StoreAuthorizationCode(ctx context.Context, clientID uint64, code string, expire time.Duration) error
	GetAuthorizationCode(ctx context.Context, clientID uint64) (code string, exist bool, err error)
}

type QuotaCache interface {
	LastStatisticsAt(ctx context.Context) (statisticsAt time.Time, err error)
	GetTotalQuota(ctx context.Context, clientID uint64) (quota uint64, err error)
}
