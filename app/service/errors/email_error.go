package errors

import "errors"

func InvalidEmailAddress(email string) error {
	return errors.New("invalid email address: " + email)
}

func InvalidEmailTemplate(templateKey string) error {
	return errors.New("invalid email template key: " + templateKey)
}

func ParseTemplateFailed(templateKey string) error {
	return errors.New("failed to parse email template: " + templateKey)
}

func SendEmailFailed(templateKey string, receiver string) error {
	return errors.New("failed to send email: " + templateKey + "to " + receiver)
}
