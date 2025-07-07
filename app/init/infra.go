package init

import (
	"github.com/alioth-center/dusk-scheduler/app/config"
	"github.com/alioth-center/dusk-scheduler/infra/apis/location"
	"github.com/alioth-center/dusk-scheduler/infra/email"
)

func initEmailSenderClient(config *config.AppConfig) {
	emailConfig := config.EmailConfig
	switch emailConfig.Provider {
	case "smtp":
		provider := emailConfig.SmtpProvider
		emailSenderClient = email.NewSmtpSenderClient(email.SmtpAuthSecret{
			Username: provider.Username,
			Password: provider.Password,
			Host:     provider.Host,
			Port:     uint16(provider.Port),
			Sender:   provider.Sender,
		})
	default:
		panic("unsupported email provider " + emailConfig.Provider)
	}
}

func initPositionLocator(config *config.AppConfig) {
	locatorConfig := config.PositionLocatorConfig
	switch locatorConfig.Provider {
	case "tencent_map":
		panic("tencent_map support not implemented yet")
	case "amap":
		panic("amap support not implemented yet")
	case "ip.sb":
		positionLocator = location.NewIpSbPositionLocator()
	default:
		positionLocator = location.NewIpSbPositionLocator()
	}
}
