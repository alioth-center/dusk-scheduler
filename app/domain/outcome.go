package domain

import "time"

const TaskNameOutcome = "outcome"

type Outcome struct {
	ID          uint64    `gorm:"column:id;type:bigint(20);autoIncrement:true;primaryKey;not null"`
	Instance    uint64    `gorm:"column:instance;type:bigint(20);not null;default:0;index:idx_instance"`
	TaskID      uint64    `gorm:"column:task_id;type:bigint(20);not null;default:0;index:idx_task"`
	Reference   string    `gorm:"column:reference;type:varchar(255);not null;default:'';unique;uniqueIndex:uk_reference"`
	StartedAt   time.Time `gorm:"column:started_at;type:timestamp;not null"`
	CompletedAt time.Time `gorm:"column:completed_at;type:timestamp;not null"`
}

func (Outcome) TableName() string { return TaskNameOutcome }

type OutcomeCompleteReason int8

const (
	OutcomeCompleteReasonUnknown   OutcomeCompleteReason = 0
	OutcomeCompleteReasonCompleted OutcomeCompleteReason = 1
	OutcomeCompleteReasonExpired   OutcomeCompleteReason = 2
	OutcomeCompleteReasonError     OutcomeCompleteReason = 3
)

func (enum OutcomeCompleteReason) String() string {
	switch enum {
	case OutcomeCompleteReasonUnknown:
		return "unknown"
	case OutcomeCompleteReasonCompleted:
		return "completed"
	case OutcomeCompleteReasonExpired:
		return "expired"
	case OutcomeCompleteReasonError:
		return "error"
	default:
		return "unknown"
	}
}
