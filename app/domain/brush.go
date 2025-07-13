package domain

import "time"

const TableNameBrush = "brush"

type Brush struct {
	ID             uint64        `gorm:"column:id;type:bigint(20);autoIncrement:true;not null;primaryKey"`
	Protocol       BrushProtocol `gorm:"column:protocol;type:tinyint(1);not null;default:1"`
	Maintainer     string        `gorm:"column:maintainer;type:varchar(30);not null;index:idx_maintainer"`
	CallURL        string        `gorm:"column:call_url;type:varchar(256);not null;default:'';index:idx_call_url"`
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

func (enum BrushProtocol) String() string {
	switch enum {
	case BrushProtocolHttp:
		return "http"
	case BrushProtocolGrpc:
		return "grpc"
	case BrushProtocolTcp:
		return "tcp"
	default:
		return "unknown"
	}
}
