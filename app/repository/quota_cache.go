package repository

import (
	"context"
	"github.com/alioth-center/dusk-scheduler/infra/cache"
	"strconv"
	"time"
)

const (
	lastStatisticsKey  = "quota_last_statistics"
	clientLastQuotaKey = "client_last_statistics_quota"
)

type quotaCache struct {
	db cache.Cache
}

func NewQuotaCache(db cache.Cache) QuotaCache {
	return &quotaCache{db: db}
}

func (cache *quotaCache) LastStatisticsAt(ctx context.Context) (statisticsAt time.Time, err error) {
	stringValue, exist, getErr := cache.db.Strings().Get(ctx, lastStatisticsKey)
	if getErr != nil {
		return time.Time{}, getErr
	}
	if len(stringValue) == 0 || !exist {
		return time.Time{}, nil
	}

	intValue, convertErr := strconv.ParseInt(stringValue, 10, 64)
	if convertErr != nil {
		return time.Time{}, convertErr
	}

	return time.Unix(intValue, 0), nil
}

func (cache *quotaCache) GetTotalQuota(ctx context.Context, clientID uint64) (quota uint64, err error) {
	stringValue, existQuota, getErr := cache.db.Hashmap().HashGet(ctx, clientLastQuotaKey, strconv.FormatUint(clientID, 10))
	if getErr != nil {
		return 0, getErr
	}
	if len(stringValue) == 0 || !existQuota {
		return 0, nil
	}

	intValue, convertErr := strconv.ParseUint(stringValue, 10, 64)
	if convertErr != nil {
		return 0, convertErr
	}

	return intValue, nil
}
