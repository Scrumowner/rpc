package provider

import (
	"fmt"
	"go.uber.org/zap"
	"hugoproxy-main/notify/config"
	"net/smtp"
)

const (
	textPlain = "text/plain"
	textHtml  = "text/html"

	TextPlain = iota + 1
	TextHtml
)

type SendIn struct {
	To    string
	From  string
	Title string
	Type  int
	Data  []byte
}

type Config struct {
	SmtpHost string
	SmtpPort string
	From     string
	Password string
}
type EmailSender struct {
	config *Config
	client smtp.Auth
	logger *zap.Logger
}

func NewEmailSender(config *config.NotifyConfig, logger *zap.Logger) *EmailSender {
	emailAuth := smtp.PlainAuth("", config.SmtpConfig.From, config.SmtpConfig.Password, config.SmtpConfig.SmtpHost)
	return &EmailSender{
		config: &Config{
			SmtpHost: config.SmtpConfig.SmtpHost,
			SmtpPort: config.SmtpConfig.SmtpPort,
			From:     config.SmtpConfig.Password,
			Password: config.SmtpConfig.SmtpHost,
		},
		client: emailAuth,
		logger: logger,
	}
}

func (e *EmailSender) Send(in *SendIn) error {
	emailBody := string(in.Data)
	var contentType string

	switch in.Type {
	case TextPlain:
		contentType = textPlain
	case TextHtml:
		contentType = textHtml
	default:
		contentType = textPlain
	}

	mime := "MIME-version: 1.0;\nContent-Type: " + contentType + "; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + in.Title + "\n"
	msg := []byte(subject + mime + "\n" + emailBody)

	if err := smtp.SendMail(fmt.Sprint("%s:%s", e.config.SmtpHost, e.config.SmtpPort),
		e.client,
		e.config.From,
		[]string{in.To},
		msg); err != nil {
		e.logger.Error("email: sent msg err", zap.Error(err))
		return err
	}

	return nil

}
