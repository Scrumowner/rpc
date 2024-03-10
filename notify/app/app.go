package app

import "notify/config"

type Application interface {
	Runner
	Bootstraper
}
type Runner interface {
	Run()
}
type Bootstraper interface {
	Bootstrap()
}
type NotifyApp struct {
	config *config.NotifyConfig
}

func NewNotifyApp() *NotifyApp {
	return &NotifyApp{}
}

func (a *NotifyApp) Bootstrap(config *config.NotifyConfig) Runner {
	return a
}

func (a *NotifyApp) Run() {

}
