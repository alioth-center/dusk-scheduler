package config

type Config interface {
	GetString(keys ...string) (value string, exist bool)
	GetInt(keys ...string) (value int, exist bool)
	GetBool(keys ...string) (value bool, exist bool)
	GetStringMap(keys ...string) map[string]any
}
