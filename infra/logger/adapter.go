package logger

import (
	"context"
	"gorm.io/gorm/logger"
	"time"
)

type Logger interface {
	DebugCtx(ctx context.Context, message string, data any)
	InfoCtx(ctx context.Context, message string, data any)
	WarnCtx(ctx context.Context, message string, data any)
	ErrorCtx(ctx context.Context, message string, err error)
}

type DatabaseLogger interface {
	LogMode(level logger.LogLevel) logger.Interface
	Info(ctx context.Context, s string, i ...interface{})
	Warn(ctx context.Context, s string, i ...interface{})
	Error(ctx context.Context, s string, i ...interface{})
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}
