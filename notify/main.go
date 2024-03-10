package main

import (
	"notify/config"
)

func main() {
	//amqpDial, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672")
	config.NewNotifyConfig()
}
