package config

type AppConfig struct {
	EmailConfig           EmailConfig           `json:"email_config" yaml:"email_config"`
	PositionLocatorConfig PositionLocatorConfig `json:"position_locator_config" yaml:"position_locator_config"`
	ClientOptions         ClientOptions         `json:"client_options" yaml:"client_options"`
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

type PositionLocatorConfig struct {
	Provider string `json:"provider" yaml:"provider"`
}

type ClientOptions struct {
	DefaultPermission     DefaultClientPermission `json:"default_permission" yaml:"default_permission"`
	ClientApiKeyPrefix    string                  `json:"client_api_key_prefix" yaml:"client_api_key_prefix"`
	AuthCodeExpireSeconds int                     `json:"auth_code_expire_seconds" yaml:"auth_code_expire_seconds"`
	DesensitizePrefix     int                     `json:"desensitize_prefix" yaml:"desensitize_prefix"`
	DesensitizeSuffix     int                     `json:"desensitize_suffix" yaml:"desensitize_suffix"`
	DesensitizeDomain     bool                    `json:"desensitize_domain" yaml:"desensitize_domain"`
}

type DefaultClientPermission struct {
	PromotionalCode string `json:"promotional_code" yaml:"promotional_code"`
	AllowBrush      bool   `json:"allow_brush" yaml:"allow_brush"`
	AllowDelay      bool   `json:"allow_delay" yaml:"allow_delay"`
	AllowHeight     int    `json:"allow_height" yaml:"allow_height"`
	AllowWidth      int    `json:"allow_width" yaml:"allow_width"`
	AllowPriority   int    `json:"allow_priority" yaml:"allow_priority"`
	DefaultQuota    int    `json:"default_quota" yaml:"default_quota"`
	LimitFrequency  int    `json:"limit_frequency" yaml:"limit_frequency"`
	LimitDuration   int    `json:"limit_duration" yaml:"limit_duration"`
}
