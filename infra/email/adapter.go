package email

import (
	"bytes"
	"context"
)

type SenderClient interface {
	SendHtmlEmail(ctx context.Context, receiver string, subject string, content *bytes.Buffer) error
	SendHtmlEmailBatch(ctx context.Context, receivers []string, subject string, content *bytes.Buffer) error
}
