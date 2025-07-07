package email

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jordan-wright/email"
	"mime"
	"net/smtp"
	"net/textproto"
)

type SmtpAuthSecret struct {
	Username string
	Password string
	Host     string
	Port     uint16
	Sender   string
}

type smtpSender struct {
	secret SmtpAuthSecret
	auth   smtp.Auth
}

func NewSmtpSenderClient(credential SmtpAuthSecret) SenderClient {
	client := smtpSender{
		secret: credential,
		auth:   smtp.PlainAuth("", credential.Username, credential.Password, credential.Host),
	}

	return &client
}

func (s *smtpSender) SendHtmlEmail(ctx context.Context, receiver string, subject string, content *bytes.Buffer) error {
	return s.SendHtmlEmailBatch(ctx, []string{receiver}, subject, content)
}

func (s *smtpSender) SendHtmlEmailBatch(_ context.Context, receivers []string, subject string, content *bytes.Buffer) error {
	emailItem := email.Email{
		From:    s.secret.Sender,
		To:      receivers,
		Subject: mime.QEncoding.Encode("utf-8", subject),
		HTML:    content.Bytes(),
		Headers: textproto.MIMEHeader{},
	}
	emailItem.Headers.Set("Content-Type", "text/html; charset=utf-8")
	emailItem.Headers.Set("Content-Transfer-Encoding", "quoted-printable")
	hostAddress := fmt.Sprintf("%s:%d", s.secret.Host, s.secret.Port)

	return emailItem.Send(hostAddress, s.auth)
}
