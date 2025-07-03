package cache

import "time"

type Cache interface {
	GetString(key string) (value string, exist bool)
	GetInt(key string) (value int, exist bool)
	GetStruct(key string, receiver any) (err error, exist bool)
	SetString(key string, value string)
	SetInt(key string, value int)
	SetStruct(key string, sender any) (err error)
	Increase(key string, delta int) (result int)
	Expire(key string, duration time.Duration) (err error)
	SetStringWithExpire(key string, value string, expireDuration time.Duration)
	SetIntWithExpire(key string, value int, expireDuration time.Duration)
	SetStructWithExpire(key string, sender any, expireDuration time.Duration) (err error)
	GetLock(key string) (success bool)
	UnLock(key string)
}
