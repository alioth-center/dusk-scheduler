package domain

import (
	"encoding/json"
	"time"
)

const TableNameTask = "task"

type Task struct {
	ID            uint64            `gorm:"column:id;type:bigint(20);autoIncrement:true;primaryKey;not null"`
	Submitter     uint64            `gorm:"column:submitter;type:bigint(20);not null;default:0;index:idx_submitter"`
	ContentHash   string            `gorm:"column:content_hash;type:varchar(64);not null;default:'';index:idx_content"`
	Width         uint              `gorm:"column:width;type:int(10);not null;default:0"`
	Height        uint              `gorm:"column:height;type:int(10);not null;default:0"`
	Type          TaskType          `gorm:"column:type;type:tinyint(1);not null;default:0"`
	Priority      TaskPriority      `gorm:"column:priority;type:tinyint(1);not null;default:0"`
	Format        TaskFormat        `gorm:"column:format;type:tinyint(1);not null;default:0"`
	ExtraOptions  json.RawMessage   `gorm:"column:extra_options;type:json;default:'{}'"`
	CreatedAt     time.Time         `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	ScheduledAt   time.Time         `gorm:"column:scheduled_at;type:timestamp;not null"`
	CompletedAt   time.Time         `gorm:"column:completed_at;type:timestamp;not null"`
	ArchivedAt    time.Time         `gorm:"column:archived_at;type:timestamp;not null"`
	ArchiveReason TaskArchiveReason `gorm:"column:archive_reason;type:tinyint(1);not null"`
}

type TaskType int8

const (
	TaskTypePainter TaskType = 0
	TaskTypeBrush   TaskType = 1
)

type TaskPriority int8

const (
	TaskPriorityLow    TaskPriority = 0
	TaskPriorityMedium TaskPriority = 1
	TaskPriorityHigh   TaskPriority = 2
)

type TaskFormat int8

const (
	TaskFormatRawImage      TaskFormat = 0
	TaskFormatImageURL      TaskFormat = 1
	TaskFormatBase64Encoded TaskFormat = 2
)

type TaskArchiveReason int8

const (
	TaskArchiveReasonUnarchived   TaskArchiveReason = 0
	TaskArchiveReasonAcknowledged TaskArchiveReason = 1
	TaskArchiveReasonCancelled    TaskArchiveReason = 2
	TaskArchiveReasonLimited      TaskArchiveReason = 3
	TaskArchiveReasonExpired      TaskArchiveReason = 4
	TaskArchiveReasonFailed       TaskArchiveReason = 5
)

func (Task) TableName() string { return TableNameTask }
