package config

type AppConfig struct {
	EmailConfig EmailConfig `json:"email_config" yaml:"email_config"`
}

type DatabaseConfig struct{}

type CacheConfig struct{}

type LoggerConfig struct{}

type EmailConfig struct {
	AllowedDomains []string                           `json:"allowed_domains" yaml:"allowed_domains"`
	MailTemplates  map[string]EmailConfigTemplateItem `json:"mail_templates" yaml:"mail_templates"`
	Provider       string                             `json:"provider" yaml:"provider"`
	SmtpProvider   EmailConfigSmtpProvider            `json:"smtp_provider" yaml:"smtp_provider"`
}

type EmailConfigTemplateItem struct {
	Subject string `json:"subject" yaml:"subject"`
	Text    string `json:"text" yaml:"text"`
}

type EmailConfigSmtpProvider struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Sender   string `json:"sender" yaml:"sender"`
}
