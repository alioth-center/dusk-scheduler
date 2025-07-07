package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/alioth-center/dusk-scheduler/app/config"
	"github.com/alioth-center/dusk-scheduler/app/service/errors"
	"github.com/alioth-center/dusk-scheduler/infra/email"
	"github.com/alioth-center/dusk-scheduler/infra/logger"
	"github.com/alioth-center/dusk-scheduler/infra/utils"
	"html/template"
)

const (
	EmailTemplateKeyRegisterClient = "register_client"
)

type emailService struct {
	sysLogger logger.Logger
	appConfig *config.AppConfig

	emailClient            email.SenderClient
	allowedDomainMapping   map[string]struct{}
	templateSubjectMapping map[string]string
	templateTextMapping    map[string]*template.Template
}

func NewEmailService(
	emailClient email.SenderClient,
	sysLogger logger.Logger,
	appConfig *config.AppConfig,
) EmailService {
	svc := emailService{
		sysLogger:   sysLogger,
		appConfig:   appConfig,
		emailClient: emailClient,

		allowedDomainMapping:   make(map[string]struct{}),
		templateSubjectMapping: make(map[string]string),
		templateTextMapping:    make(map[string]*template.Template),
	}

	for _, domain := range svc.appConfig.EmailConfig.AllowedDomains {
		svc.allowedDomainMapping[domain] = struct{}{}
	}
	for key, content := range svc.appConfig.EmailConfig.MailTemplates {
		parsedTemplate, parseErr := template.New(key).Parse(content.Text)
		if parseErr != nil {
			panic(fmt.Sprintf("failed to parse email template %s: %v", key, parseErr))
		}

		svc.templateTextMapping[key] = parsedTemplate
		svc.templateSubjectMapping[key] = content.Subject
	}

	return &svc
}

func (srv *emailService) ValidateEmailAddress(ctx context.Context, emailAddress string) (err error) {
	if valid, allowed := utils.ValidateEmailAddress(emailAddress, srv.allowedDomainMapping); !valid {
		return errors.InvalidEmailAddress(emailAddress)
	} else if !allowed {
		srv.sysLogger.InfoCtx(ctx, "domain not allowed", map[string]any{"email_address": emailAddress})

		return errors.InvalidEmailAddress(emailAddress)
	}

	return nil
}

func (srv *emailService) SendEmail(ctx context.Context, receiver string, templateKey string, args map[string]any) (err error) {
	templateContent, existTemplate := srv.templateTextMapping[templateKey]
	if !existTemplate {
		srv.sysLogger.InfoCtx(ctx, "template does not exist", map[string]any{"template_key": templateKey})

		return errors.InvalidEmailTemplate(templateKey)
	}
	templateSubject, existSubject := srv.templateSubjectMapping[templateKey]
	if !existSubject {
		srv.sysLogger.InfoCtx(ctx, "subject does not exist", map[string]any{"template_key": templateKey})

		return errors.InvalidEmailTemplate(templateKey)
	}

	content := &bytes.Buffer{}
	if parseErr := templateContent.Execute(content, args); parseErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to parse email template: %s, args: %v", templateKey, args), parseErr)

		return errors.ParseTemplateFailed(templateKey)
	}

	if sendErr := srv.emailClient.SendHtmlEmail(ctx, receiver, templateSubject, content); sendErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to send email template: %s, receiver: %s", templateKey, receiver), sendErr)

		return errors.SendEmailFailed(templateKey, receiver)
	}

	return nil
}
