package logger

import "context"

type Logger interface {
	DebugCtx(ctx context.Context, message string, data any)
	InfoCtx(ctx context.Context, message string, data any)
	WarnCtx(ctx context.Context, message string, data any)
	ErrorCtx(ctx context.Context, err error, data any)
}
