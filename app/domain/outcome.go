package domain

import "time"

const TaskNameOutcome = "outcome"

type Outcome struct {
	ID          uint64    `gorm:"column:id;type:bigint(20);autoIncrement;primaryKey;not null"`
	Instance    uint64    `gorm:"column:instance;type:bigint(20);not null;default:0;index:idx_instance"`
	TaskID      uint64    `gorm:"column:task_id;type:bigint(20);not null;default:0;index:idx_task"`
	Reference   string    `gorm:"column:reference;type:varchar(255);not null;default:'';index:idx_reference"`
	StartedAt   time.Time `gorm:"column:started_at;type:timestamp;not null"`
	CompletedAt time.Time `gorm:"column:completed_at;type:timestamp;not null"`
}

func (Outcome) TableName() string { return TaskNameOutcome }
