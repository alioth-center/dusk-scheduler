package domain

const TableNameClient = "client"

type Client struct {
	ID              uint64       `gorm:"column:id;type:bigint(20);autoIncrement;primaryKey;not null"`
	Maintainer      string       `gorm:"column:maintainer;type:varchar(128);not null;default:'';uniqueIndex:uk_client"`
	Region          string       `gorm:"column:region;type:varchar(128);not null;default:''"`
	Authorized      bool         `gorm:"column:authorized;type:tinyint(1);not null;default:0"`
	ApiKey          string       `gorm:"column:api_key;type:varchar(128);not null;default:'';uniqueIndex:uk_client"`
	AllowBrush      bool         `gorm:"column:allow_brush;type:tinyint(1);not null;default:0"`
	AllowDelay      bool         `gorm:"column:allow_delay;type:tinyint(1);not null;default:0"`
	AllowHeight     uint32       `gorm:"column:allow_height;type:int(10);not null;default:0"`
	AllowWidth      uint32       `gorm:"column:allow_width;type:int(10);not null;default:0"`
	AllowPriority   TaskPriority `gorm:"column:allow_priority;type:tinyint(1);not null;default:0"`
	Quota           uint32       `gorm:"column:quota;type:int(10);not null;default:10000"`
	LimitRenderSize uint32       `gorm:"column:limit_render_size;type:int(10);not null;default:1024"`
	LimitFrequency  uint32       `gorm:"column:limit_frequency;type:int(10);not null;default:10"`
	LimitDuration   uint32       `gorm:"column:limit_duration;type:int(10);not null;default:60"`
}

func (Client) TableName() string { return TableNameClient }
