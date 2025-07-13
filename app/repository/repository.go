package repository

import (
	"bytes"
	"context"
	"github.com/alioth-center/dusk-scheduler/app/domain"
	"time"
)

type ClientDao interface {
	CreateClient(ctx context.Context, client *domain.Client) (clientID uint64, err error)
	GetClientByID(ctx context.Context, clientID uint64) (client *domain.Client, exist bool, err error)
}

type TaskDao interface {
	CreateTask(ctx context.Context, task *domain.Task) (taskID uint64, err error)
	GetTaskByID(ctx context.Context, taskID uint64) (task *domain.Task, exist bool, err error)
	GetTaskListByClientID(ctx context.Context, clientID uint64, statusFilter []string, offsetTaskID uint64, pageLimit uint32, desc bool) (tasks []*domain.Task, err error)
	UpdateTaskAsCompleted(ctx context.Context, taskID uint64) error
	UpdateTaskAsArchived(ctx context.Context, taskID uint64, reason domain.TaskArchiveReason) error
	StatisticsClientQuotaUsage(ctx context.Context, clientID uint64, startTime time.Time) (usage uint64, err error)
}

type PromotionalDao interface {
	GetPromotionalByCode(ctx context.Context, code string) (promotional *domain.Promotional, exist bool, err error)
}

type OutcomeDao interface {
	CreateOutcome(ctx context.Context, outcome *domain.Outcome) (outcomeID uint64, err error)
	GetOutcomeByReference(ctx context.Context, outcomeReference string) (outcome *domain.Outcome, exist bool, err error)
	GetOutcomeByTaskID(ctx context.Context, outcomeID uint64) (outcome *domain.Outcome, exist bool, err error)
}

type PainterDao interface {
	CreatePainter(ctx context.Context, painter *domain.Painter) (painterID uint64, err error)
	UpdatePainterName(ctx context.Context, painterID uint64, name string) error
	GetPainterByID(ctx context.Context, painterID uint64) (painter *domain.Painter, exist bool, err error)
	GetPainterByName(ctx context.Context, name string) (painter *domain.Painter, exist bool, err error)
	UpdatePainterAsConnected(ctx context.Context, painterID uint64) error
	UpdatePainterAsDisconnected(ctx context.Context, painterID uint64) error
}

type StorageDao interface {
	GetStorageByID(ctx context.Context, storageID uint64) (storage *domain.Storage, exist bool, err error)
	GetStorageByName(ctx context.Context, name string) (storage *domain.Storage, exist bool, err error)
}

type AuthorizationCache interface {
	StoreAuthorizationCode(ctx context.Context, clientID uint64, code string, expire time.Duration) error
	GetAuthorizationCode(ctx context.Context, clientID uint64) (code string, exist bool, err error)
}

type QuotaCache interface {
	LastStatisticsAt(ctx context.Context) (statisticsAt time.Time, err error)
	GetTotalQuota(ctx context.Context, clientID uint64) (quota uint64, err error)
}

type TaskContentCache interface {
	StoreTaskContent(ctx context.Context, taskID uint64, content *bytes.Buffer) (err error)
	DeleteTaskContent(ctx context.Context, taskID uint64) (err error)
}

type PainterHeartbeatCache interface {
	GetHeartbeatTime(ctx context.Context, name string) (lastHeartbeatAt time.Time, err error)
	UpdateHeartbeatTime(ctx context.Context, name string) (err error)
	DeleteHeartbeatTime(ctx context.Context, name string) (err error)
}
