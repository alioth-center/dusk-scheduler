package storage

type Storage interface {
	LoadFile(key string) (content []byte, err error)
	StoreFile(key string, content []byte) (err error)
	DeleteFile(key string) (err error)
}
