package repository

import (
	"context"
	"fmt"
	"github.com/alioth-center/dusk-scheduler/infra/cache"
	"time"
)

const (
	authorizationCodeKey = "authorization_code:%d"
)

type authorizationCache struct {
	db cache.Cache
}

func NewAuthorizationCache(db cache.Cache) AuthorizationCache {
	return &authorizationCache{db: db}
}

func (cache *authorizationCache) StoreAuthorizationCode(ctx context.Context, clientID uint64, code string, expire time.Duration) error {
	key := fmt.Sprintf(authorizationCodeKey, clientID)
	if setErr := cache.db.Strings().SetEx(ctx, key, code, expire); setErr != nil {
		return setErr
	}

	return nil
}

func (cache *authorizationCache) GetAuthorizationCode(ctx context.Context, clientID uint64) (code string, exist bool, err error) {
	key := fmt.Sprintf(authorizationCodeKey, clientID)
	result, existCode, getErr := cache.db.Strings().Get(ctx, key)
	if getErr != nil {
		return "", false, getErr
	}
	if !existCode {
		return "", false, nil
	}

	return result, true, nil
}
