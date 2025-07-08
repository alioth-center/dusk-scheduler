package cache

import (
	"bytes"
	"context"
	"time"
)

type Cache interface {
	Strings() StringCache
	Hashmap() HashmapCache
	SortedSet() SortedSetCache
}

type CachedKeys interface {
	Delete(ctx context.Context, key string) (err error)
	Exists(ctx context.Context, key string) (exist bool, err error)
	Expire(ctx context.Context, key string, duration time.Duration) (err error)
	ExpireAt(ctx context.Context, key string, expireAt time.Time) (err error)
	TTL(ctx context.Context, key string) (ttl time.Duration, err error)
	Rename(ctx context.Context, key, newName string) (err error)
	RenameNX(ctx context.Context, key, newName string) (changed bool, err error)
}

type StringCache interface {
	CachedKeys
	Set(ctx context.Context, key, value string) (err error)
	SetEx(ctx context.Context, key, value string, ttl time.Duration) (err error)
	SetNx(ctx context.Context, key, value string, ttl time.Duration) (err error)
	SetBytes(ctx context.Context, key string, value *bytes.Buffer) (err error)
	SetBytesNx(ctx context.Context, key string, value *bytes.Buffer, ttl time.Duration) (err error)
	SetBytesEx(ctx context.Context, key string, value *bytes.Buffer, ttl time.Duration) (err error)
	Get(ctx context.Context, key string) (value string, exist bool, err error)
	GetBytes(ctx context.Context, key string) (value *bytes.Buffer, exist bool, err error)
	Swap(ctx context.Context, key string, newValue string) (old string, exist bool, err error)
	SwapBytes(ctx context.Context, key string, newValue *bytes.Buffer) (old *bytes.Buffer, exist bool, err error)
	IncrBy(ctx context.Context, key string, delta uint64) (newValue uint64, exist bool, err error)
	DecrBy(ctx context.Context, key string, delta uint64) (newValue uint64, exist bool, err error)
}

type HashmapCache interface {
	CachedKeys
	HashDelete(ctx context.Context, key, field string) (err error)
	HashExists(ctx context.Context, key, field string) (exist bool, err error)
	HashSet(ctx context.Context, key, field string, value string) (err error)
	HashSetBytes(ctx context.Context, key, field string, value *bytes.Buffer) (err error)
	HashSetNx(ctx context.Context, key, field string, value string) (err error)
	HashSetBytesNx(ctx context.Context, key, field string, value *bytes.Buffer) (err error)
	HashGet(ctx context.Context, key, field string) (value string, exist bool, err error)
	HashGetBytes(ctx context.Context, key, field string) (value *bytes.Buffer, exist bool, err error)
	HashIncrBy(ctx context.Context, key, field string, delta uint64) (newValue uint64, exist bool, err error)
	HashDecrBy(ctx context.Context, key, field string, delta uint64) (newValue uint64, exist bool, err error)
}

type SortedSetCache interface {
	CachedKeys
	SortedSetAdd(ctx context.Context, key string, value string, score uint64) (err error)
	SortedSetAddBytes(ctx context.Context, key string, value *bytes.Buffer, score uint64) (err error)
	SortedSetRemove(ctx context.Context, key string) (err error)
	SortedSetRemoveRange(ctx context.Context, key string, minScore, maxScore uint64) (err error)
	SortedSetCountRange(ctx context.Context, key string, minScore, maxScore uint64) (count uint64, err error)
	SortedSetPopN(ctx context.Context, key string, count uint64) (values []string, err error)
	SortedSetPopBytesN(ctx context.Context, key string, count uint64) (values []*bytes.Buffer, err error)
}
