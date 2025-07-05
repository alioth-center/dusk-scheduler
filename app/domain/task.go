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
	Width         uint32            `gorm:"column:width;type:int(10);not null;default:0"`
	Height        uint32            `gorm:"column:height;type:int(10);not null;default:0"`
	Type          TaskType          `gorm:"column:type;type:tinyint(1);not null;default:0"`
	Priority      TaskPriority      `gorm:"column:priority;type:tinyint(1);not null;default:0"`
	Format        TaskFormat        `gorm:"column:format;type:tinyint(1);not null;default:0"`
	DelayRender   uint8             `gorm:"column:delay_render;type:tinyint(1);not null;default:0"`
	ExtraOptions  json.RawMessage   `gorm:"column:extra_options;type:json;default:'{}'"`
	Instance      uint64            `gorm:"column:instance;type:bigint(20);not null;default:0;index:idx_instance"`
	CreatedAt     time.Time         `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;index:idx_status"`
	ScheduledAt   time.Time         `gorm:"column:scheduled_at;type:timestamp;not null;index:idx_status"`
	CompletedAt   time.Time         `gorm:"column:completed_at;type:timestamp;not null;index:idx_status"`
	ArchivedAt    time.Time         `gorm:"column:archived_at;type:timestamp;not null;index:idx_status"`
	QuotaUsage    uint32            `gorm:"column:quota_usage;type:int(10);not null;default:0"`
	ArchiveReason TaskArchiveReason `gorm:"column:archive_reason;type:tinyint(1);not null"`
}

func (Task) TableName() string { return TableNameTask }

func (t Task) Status() TaskStatus {
	switch {
	case !t.ArchivedAt.IsZero():
		return TaskStatusArchived
	case !t.CompletedAt.IsZero():
		return TaskStatusCompleted
	case !t.ScheduledAt.IsZero():
		return TaskStatusProcessing
	default:
		return TaskStatusGenerated
	}
}

type TaskType int8

const (
	TaskTypePainter TaskType = 0
	TaskTypeBrush   TaskType = 1
)

func (enum TaskType) String() string {
	switch enum {
	case TaskTypePainter:
		return "painter"
	case TaskTypeBrush:
		return "brush"
	default:
		return "unknown"
	}
}

type TaskPriority int8

const (
	TaskPriorityLow    TaskPriority = 0
	TaskPriorityMedium TaskPriority = 1
	TaskPriorityHigh   TaskPriority = 2
)

func (enum TaskPriority) String() string {
	switch enum {
	case TaskPriorityLow:
		return "low"
	case TaskPriorityMedium:
		return "medium"
	case TaskPriorityHigh:
		return "high"
	default:
		return "unknown"
	}
}

type TaskFormat int8

const (
	TaskFormatRawImage      TaskFormat = 0
	TaskFormatImageURL      TaskFormat = 1
	TaskFormatBase64Encoded TaskFormat = 2
)

func (enum TaskFormat) String() string {
	switch enum {
	case TaskFormatRawImage:
		return "raw_image"
	case TaskFormatImageURL:
		return "image_url"
	case TaskFormatBase64Encoded:
		return "base64_encoded"
	default:
		return "unknown"
	}
}

type TaskArchiveReason int8

const (
	TaskArchiveReasonUnarchived   TaskArchiveReason = 0
	TaskArchiveReasonAcknowledged TaskArchiveReason = 1
	TaskArchiveReasonCancelled    TaskArchiveReason = 2
	TaskArchiveReasonLimited      TaskArchiveReason = 3
	TaskArchiveReasonExpired      TaskArchiveReason = 4
	TaskArchiveReasonFailed       TaskArchiveReason = 5
)

func (enum TaskArchiveReason) String() string {
	switch enum {
	case TaskArchiveReasonUnarchived:
		return "unarchived"
	case TaskArchiveReasonAcknowledged:
		return "acknowledged"
	case TaskArchiveReasonCancelled:
		return "cancelled"
	case TaskArchiveReasonLimited:
		return "limited"
	case TaskArchiveReasonExpired:
		return "expired"
	case TaskArchiveReasonFailed:
		return "failed"
	default:
		return "unknown"
	}
}

type TaskStatus = string

const (
	TaskStatusGenerated  TaskStatus = "generated"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusArchived   TaskStatus = "archived"
)
