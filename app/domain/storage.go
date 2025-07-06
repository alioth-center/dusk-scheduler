package domain

import (
	"encoding/json"
	"time"
)

const TableNameStorage = "storage"

type Storage struct {
	ID         uint64          `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true;not null"`
	Protocol   StorageProtocol `gorm:"column:protocol;type:tinyint(1);not null;default:0"`
	PolicyName string          `gorm:"column:policy_name;type:varchar(50);not null;unique;uniqueIndex:uk_name"`
	Options    json.RawMessage `gorm:"column:options;type:json;not null;default:'{}'"`
	CreatedAt  time.Time       `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time       `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (Storage) TableName() string { return TableNameStorage }

type StorageProtocol int8

const (
	StorageProtocolUnknown   StorageProtocol = 0
	StorageProtocolLocalFile StorageProtocol = 1
	StorageProtocolS3        StorageProtocol = 2
	StorageProtocolFtp       StorageProtocol = 3
)

func (enum StorageProtocol) String() string {
	switch enum {
	case StorageProtocolLocalFile:
		return "local_file"
	case StorageProtocolS3:
		return "s3"
	case StorageProtocolFtp:
		return "ftp"
	default:
		return "unknown"
	}
}

type StorageOptionsS3 struct {
	Protocol  string `json:"protocol"`
	Endpoint  string `json:"endpoint"`
	Bucket    string `json:"bucket"`
	Prefix    string `json:"prefix"`
	Region    string `json:"region"`
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
}

type StorageOptionsFtp struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	RemotePath string `json:"remote_path"`
}

type StorageOptionsLocalFile struct {
	Path string `json:"path"`
}
