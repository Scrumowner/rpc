package controllers

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	responder "proxy/internal/infra/responder"
	"proxy/internal/modules/controller"
)

type Controllers struct {
	UserController   *controller.User
	AuthController   *controller.Auth
	SwagController   *controller.Swagger
	SearchController *controller.Search
}

func NewControllers(logger *zap.SugaredLogger, user grpc.ClientConnInterface, auth grpc.ClientConnInterface, geo grpc.ClientConnInterface) *Controllers {
	respond := responder.NewResponder(logger)

	return &Controllers{
		UserController:   controller.NewUserController(respond, logger, user),
		AuthController:   controller.NewAuthController(logger, respond, auth),
		SwagController:   controller.NewSwaggerer(logger, respond),
		SearchController: controller.NewSearchController(respond, logger, geo),
	}
}
