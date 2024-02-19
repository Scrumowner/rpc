package controllers

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"hugoproxy-main/proxy/controller"
	"hugoproxy-main/proxy/responder"
	"net/http"
)

type Controllers struct {
	AuthController   *controller.Auth
	UserController   *controller.User
	SwagController   *controller.Swagger
	SearchController *controller.Search
}

func NewControllers(logger *zap.SugaredLogger, client http.Client, geo grpc.ClientConnInterface, auth grpc.ClientConnInterface, user grpc.ClientConnInterface) *Controllers {
	responder := responder.NewResponder(logger)
	return &Controllers{
		AuthController:   controller.NewAuthController(logger, client, responder, auth),
		SwagController:   controller.NewSwaggerer(logger, client, responder),
		UserController:   controller.NewUserController(responder, logger, user),
		SearchController: controller.NewSearchController(responder, logger, geo),
	}
}
