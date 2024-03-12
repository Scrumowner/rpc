package provider

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"log"
	"net/http"
)

const Sender = "http://81.163.28.166:8080/api/sms/send"

type SendSmsIn struct {
	To      string
	Message string
}

type PhoneSender struct {
	client *http.Client
	logger *zap.Logger
}

func NewPhoneSender() *PhoneSender {
	return &PhoneSender{
		client: http.DefaultClient,
	}
}

func (p *PhoneSender) SendSms(in *SendSmsIn) error {
	buf := bytes.NewBuffer(make([]byte, 0))
	err := json.NewEncoder(buf).Encode(in)
	if err != nil {
		log.Println("CAN'T MARSHALL SendSmsIn sturct TO JSON")
		return nil
	}
	req, err := http.NewRequest("POST", Sender, buf)
	if err != nil {
		log.Println(err, "CAN'T PREPARE REQUEST TO SEND INTO SMS SEN API")
	}
	_, err = p.client.Do(req)
	if err != nil {
		log.Println(err, "SMS SEND API DOSE NOT WORK/")
		return nil
	}
	return nil
}
