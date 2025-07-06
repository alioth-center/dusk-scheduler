package domain

import "time"

const TableNameBrush = "brush"

type Brush struct {
	ID             uint64        `gorm:"column:id;type:bigint(20);autoIncrement:true;not null;primaryKey"`
	Name           string        `gorm:"column:name;type:varchar(30);not null;uniqueIndex:uk_name"`
	Protocol       BrushProtocol `gorm:"column:protocol;type:tinyint(1);not null;default:1"`
	Maintainer     string        `gorm:"column:maintainer;type:varchar(30);not null;index:idx_maintainer"`
	Secret         string        `gorm:"column:secret;type:varchar(128);not null;uniqueIndex:uk_secret"`
	ConnectTime    uint32        `gorm:"column:connect_time;type:int(10);not null;default:0"`
	RegisterAt     time.Time     `gorm:"column:register_at;type:timestamp;not null"`
	ConnectedAt    time.Time     `gorm:"column:connected_at;type:timestamp;not null"`
	DisconnectedAt time.Time     `gorm:"column:disconnected_at;type:timestamp;not null"`
}

func (Brush) TableName() string { return TableNameBrush }

type BrushProtocol int8

const (
	BrushProtocolHttp BrushProtocol = 0
	BrushProtocolGrpc BrushProtocol = 1
	BrushProtocolTcp  BrushProtocol = 2
)
