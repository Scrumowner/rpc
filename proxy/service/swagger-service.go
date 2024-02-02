package service

import (
	"hugoproxy-main/proxy/static"
)

type SwaggerServiceer interface {
	GetSwaggerHtml() string
	GetSwaggerJson() string
}

type SwaggerService struct {
}

func NewSwaggerService() SwaggerServiceer {
	return &SwaggerService{}
}

func (swagger *SwaggerService) GetSwaggerHtml() string {
	return static.SwaggerTemplate
}
func (swagger *SwaggerService) GetSwaggerJson() string {
	return static.SwagJson
}
