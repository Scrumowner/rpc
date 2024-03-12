package app

import (
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"log"
	"notify/config"
	kafka2 "notify/internal/infra/kafka"
	"notify/internal/infra/queue"
	"notify/internal/infra/rabbitmq"
	"notify/internal/modules/service"
)

const RabbitTopic = "GeoRateLimit"

type Application interface {
	Runner
}
type Runner interface {
	Run()
}

type NotifyApp struct {
	config  *config.NotifyConfig
	service *service.Service
	broker  queue.MessageQueuer
}

func NewNotifyApp(cfg *config.NotifyConfig, logger *zap.Logger, secret string, engine string) (*NotifyApp, error) {
	if engine == "rabbitmq" {
		url := fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.RabbitConfig.User, cfg.RabbitConfig.Password, cfg.RabbitConfig.Host, cfg.RabbitConfig.Port)
		amqpDial, err := amqp.Dial(url)
		if err != nil {
			log.Println("Err when dial with amqp")
		}
		rabbit, err := rabbitmq.NewRabbitMQ(amqpDial)
		if err != nil {
			log.Println("Err when try create rabbit object")
		}
		return &NotifyApp{
			config:  cfg,
			service: service.NewNotifySerivce(cfg, logger, secret),
			broker:  rabbit,
		}, nil
	}
	addr := fmt.Sprintf("%s:%s", cfg.KafkaConfig.Host, cfg.KafkaConfig.Host)
	kafka, err := kafka2.NewKafkaMQ(addr, RabbitTopic)
	if err != nil {
		log.Println("Err when dial to kafka")
	}
	return &NotifyApp{
		config:  cfg,
		service: service.NewNotifySerivce(cfg, logger, secret),
		broker:  kafka,
	}, nil
}

func (a *NotifyApp) Run() error {
	messages, err := a.broker.Subscribe(RabbitTopic)
	if err != nil {
		log.Println(err)
		return err
	}
	for {
		select {
		case notif := <-messages:
			err = a.broker.Ack(&notif)
			if err != nil {
				log.Println("ERROR WHEN ASC RECIBVING MESSAGE")
			}
			data := notif.Data
			err = a.service.Notify(data)
			if err != nil {
				log.Println(err)
			}

		default:
			continue
		}

	}

}
