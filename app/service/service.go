package service

import (
	"bytes"
	"context"
	"github.com/alioth-center/dusk-scheduler/app/domain"
	"github.com/alioth-center/dusk-scheduler/infra/apis/location"
	"net/url"
	"time"
)

type EmailService interface {
	ValidateEmailAddress(ctx context.Context, emailAddress string) (err error)
	SendEmail(ctx context.Context, receiver string, templateKey string, args map[string]any) (err error)
}

type LocationService interface {
	DetectIPLocation(ctx context.Context, ip string) (address *location.Address, err error)
}

type ClientService interface {
	CreateClient(ctx context.Context, email string, promotionCode string, ip string) (client *domain.Client, err error)
	StoreAuthorizationCode(ctx context.Context, clientID uint64, authorizationCode string) (expiredAt time.Time, err error)
	AuthorizeClient(ctx context.Context, clientID uint64, emailAddress string, authorizationCode string) (authorized bool, maintainer string, apiKey string, err error)
	GetClientData(ctx context.Context, clientID uint64) (client *domain.Client, err error)
	GetClientQuotaUsage(ctx context.Context, clientID uint64) (quotaTotal, quotaUsage uint64, lastCheckTime time.Time, err error)
}

type TaskService interface {
	GetTaskByID(ctx context.Context, taskID uint64) (task *domain.Task, exist bool, err error)
	GetCompletedTasksByClientID(ctx context.Context, clientID uint64, statusFilter []string, offsetTaskID uint64) (tasks []*domain.Task, hasMore bool, err error)
	CompleteTask(ctx context.Context, taskID uint64) (err error)
	ArchiveTaskByOutcomeReference(ctx context.Context, outcomeReference string, archiveReason domain.TaskArchiveReason) (exist bool, err error)
	FlushPainterScheduler(ctx context.Context, painterName string)
	GetScheduledTaskListByPainterName(ctx context.Context, painterName string) (list []*domain.Task, content []string, err error)
	GetNextScheduledTaskListByPainterName(ctx context.Context, painterName string) (list []*domain.Task, content []string, err error)
}

type OutcomeService interface {
	CreateOutcome(ctx context.Context, painterName string, taskID uint64, reference string, createdAt, completedAt time.Time) (err error)
	GetOutcomeContentByReference(ctx context.Context, reference string) (outcome *domain.Outcome, exist bool, err error)
	GetOutcomeContent(ctx context.Context, reference string) (content *bytes.Buffer, err error)
	GetOutcomeURL(ctx context.Context, reference string) (content *url.URL, err error)
}

type PainterService interface {
	CreatePainter(ctx context.Context, maintainer string, slot int) (painter *domain.Painter, policy *domain.Storage, err error)
	ReconnectPainter(ctx context.Context, name string) (heartbeat bool, err error)
	DisconnectPainter(ctx context.Context, name string) (err error)
	GetPainterByID(ctx context.Context, painterID uint64) (painter *domain.Painter, exist bool, err error)
	GetPainterByName(ctx context.Context, name string) (painter *domain.Painter, exist bool, err error)
	ValidatePainterByName(ctx context.Context, name string, secret string) (validated bool, err error)
}

type BrushService interface {
	CreateBrush(ctx context.Context, maintainer string, slot int, protocol domain.BrushProtocol) (brush *domain.Brush, err error)
}
