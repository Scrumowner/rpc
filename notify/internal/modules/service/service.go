package service

import (
	"go.uber.org/zap"
	"notify/config"
	"notify/internal/infra/decoder"
	"notify/internal/provider"
)

type User struct {
	Email string
	Phone string
}
type Service struct {
	email   *provider.EmailSender
	phone   *provider.PhoneSender
	decoder *decoder.Decoder
}

func NewNotifySerivce(config *config.NotifyConfig, logger *zap.Logger, secret string) *Service {
	return &Service{
		email:   provider.NewEmailSender(config, logger),
		phone:   provider.NewPhoneSender(),
		decoder: decoder.NewDecoder(secret),
	}
}

func (s *Service) Notify(token []byte) error {
	userEmail, userPhone := s.decoder.Decdoe(string(token))
	inEmail := &provider.SendIn{
		To:    userEmail,
		Title: "To many request to api service",
		Data:  []byte("To many request"),
	}
	err := s.email.Send(inEmail)
	if err != nil {
		return err
	}
	inSms := &provider.SendSmsIn{
		To:      userPhone,
		Message: "To many request",
	}
	err = s.phone.SendSms(inSms)
	if err != nil {
		return err
	}
	return nil
}
