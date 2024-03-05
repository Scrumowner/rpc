package service

import "hugoproxy-main/notify/internal/provider"

type Service struct {
	email *provider.EmailSender
}
