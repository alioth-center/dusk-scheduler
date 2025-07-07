package config

type Config interface {
	ParseAppConfig(source, namespace string, receiver any) error
}
