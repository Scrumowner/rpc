package controllers

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"hugoproxy-main/proxy/controller"
	"hugoproxy-main/proxy/responder"
	"net/http"
	"net/rpc"
)

type Controllers struct {
	SearchController        controller.Searcher
	AuthController          controller.Auther
	SwagController          controller.Swaggerer
	SearchControllerJsonRpc controller.Searcher
	SearchControllergRpc    *controller.SearchgRpc
}

func NewControllers(logger *zap.SugaredLogger, client http.Client, rpc *rpc.Client, cc grpc.ClientConnInterface) *Controllers {
	responder := responder.NewResponder(logger)
	return &Controllers{
		SearchController:        controller.NewSearcher(responder, rpc),
		SearchControllerJsonRpc: controller.NewSearcherJsonRpc(responder, rpc),
		AuthController:          controller.NewAuther(logger, client, responder),
		SwagController:          controller.NewSwaggerer(logger, client, responder),
		SearchControllergRpc:    controller.NewSearchgRpc(responder, cc),
	}
}
