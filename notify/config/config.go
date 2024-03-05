package config

import "os"

type NotifyConfig struct {
	Domain string
	Port   string
	*SmtpConfig
	*RabbitConfig
	*KafkaConfig
}
type SmtpConfig struct {
	SmtpHost string
	SmtpPort string
	From     string
	Password string
}

type RabbitConfig struct {
	User     string
	Password string
	Host     string
	Port     string
}
type KafkaConfig struct {
	Host string
	Port string
}

func NewNotifyConfig() *NotifyConfig {
	return &NotifyConfig{}
}

func (c *NotifyConfig) Load() {
	c.Domain = os.Getenv("NOTIFY_DOMAIN")
	c.Port = os.Getenv("NOTIFY_PORT")
	if "rabbitmq" == os.Getenv("MESSAGE_BROKER") {
		c.RabbitConfig = &RabbitConfig{
			User:     os.Getenv("RABBIT_USER"),
			Password: os.Getenv("RABBIT_PASSWORD"),
			Host:     os.Getenv("RABBIT_HOST"),
			Port:     os.Getenv("RABBIT_PORT"),
		}
	} else {
		c.KafkaConfig = &KafkaConfig{
			Host: os.Getenv("KAFKA_HOST"),
			Port: os.Getenv("KAFKA_PORT"),
		}
	}
	c.SmtpConfig = &SmtpConfig{
		SmtpHost: os.Getenv("SMTP_HOST"),
		SmtpPort: os.Getenv("SMTP_PORT"),
		From:     os.Getenv("EMAIL_FROM"),
		Password: os.Getenv("EMAIL_PASSWORD"),
	}

}
