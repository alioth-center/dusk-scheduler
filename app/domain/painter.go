package domain

import "time"

const TableNamePainter = "painter"

type Painter struct {
	ID             uint64    `gorm:"column:id;type:bigint(20);autoIncrement:true;not null;primaryKey;"`
	Name           string    `gorm:"column:name;type:varchar(30);not null;uniqueIndex:uk_name"`
	Location       string    `gorm:"column:location;type:varchar(30);not null;index:idx_location"`
	Maintainer     string    `gorm:"column:maintainer;type:varchar(30);not null;index:idx_maintainer"`
	Secret         string    `gorm:"column:secret;type:varchar(128);not null;uniqueIndex:uk_secret"`
	Slot           uint8     `gorm:"column:slot;type:tinyint(4);not null;default:1"`
	ConnectTime    uint32    `gorm:"column:connect_time;type:int(10);not null;default:0"`
	RegisterAt     time.Time `gorm:"column:register_at;type:timestamp;not null"`
	PolicyID       uint64    `gorm:"column:policy_id;type:bigint(20);not null;index:idx_policy_id"`
	ConnectedAt    time.Time `gorm:"column:connected_at;type:timestamp;not null"`
	DisconnectedAt time.Time `gorm:"column:disconnected_at;type:timestamp;not null"`
}

func (Painter) TableName() string { return TableNamePainter }
