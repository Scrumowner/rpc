package main

import (
	"go.uber.org/zap"
	"log"
	"notify/app"
	"notify/config"
)

func main() {

	cfg := config.NewNotifyConfig()
	cfg.Load()
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln("Can't create zap logger struct")
	}
	a, err := app.NewNotifyApp(cfg, logger, cfg.Secret, cfg.BrokerType)
	if err != nil {
		log.Println("NOTIFY APP CAN'T START")
	}
	err = a.Run()
	if err != nil {
		log.Println("NOTIFY APP CAN'T START ")
	}
	log.Println("NOTIFY APP IS DONE")
}
