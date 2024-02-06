package controllers

import (
	"go.uber.org/zap"
	"hugoproxy-main/proxy/controller"
	"hugoproxy-main/proxy/responder"
	"net/http"
	"net/rpc"
)

type Controllers struct {
	SearchController controller.Searcher
	AuthController   controller.Auther
	SwagController   controller.Swaggerer
}

func NewControllers(logger *zap.SugaredLogger, client http.Client, rpc *rpc.Client) *Controllers {
	responder := responder.NewResponder(logger)
	return &Controllers{
		SearchController: controller.NewSearcher(responder, rpc),
		AuthController:   controller.NewAuther(logger, client, responder),
		SwagController:   controller.NewSwaggerer(logger, client, responder),
	}
}
