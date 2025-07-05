package domain

import "time"

const TableNamePromotional = "promotional"

type Promotional struct {
	ID              uint64       `gorm:"column:id;type:bigint(20);primary_key;not null;primaryKey"`
	Code            string       `gorm:"column:code;type:varchar(50);not null;unique;uniqueIndex:idx_code"`
	CreatedAt       time.Time    `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	AllowBrush      bool         `gorm:"column:allow_brush;type:tinyint(1);not null;default:0"`
	AllowDelay      bool         `gorm:"column:allow_delay;type:tinyint(1);not null;default:0"`
	AllowHeight     uint32       `gorm:"column:allow_height;type:int(10);not null;default:0"`
	AllowWidth      uint32       `gorm:"column:allow_width;type:int(10);not null;default:0"`
	AllowPriority   TaskPriority `gorm:"column:allow_priority;type:tinyint(1);not null;default:0"`
	Quota           uint32       `gorm:"column:quota;type:int(10);not null;default:0"`
	LimitRenderSize uint32       `gorm:"column:limit_render_size;type:int(10);not null;default:0"`
	LimitFrequency  uint32       `gorm:"column:limit_frequency;type:int(10);not null;default:0"`
	LimitDuration   uint32       `gorm:"column:limit_duration;type:int(10);not null;default:0"`
	StartedAt       time.Time    `gorm:"column:started_at;type:timestamp;not null;index:idx_time"`
	ExpiredAt       time.Time    `gorm:"column:expired_at;type:timestamp;not null;index:idx_time"`
}

func (Promotional) TableName() string { return TableNamePromotional }
