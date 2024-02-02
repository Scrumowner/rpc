package controllers

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"hugoproxy-main/proxy/controller"
	"hugoproxy-main/proxy/responder"
	"net/http"
)

type Controllers struct {
	SearchController controller.Searcher
	AuthController   controller.Auther
	SwagController   controller.Swaggerer
}

func NewControllers(logger *zap.SugaredLogger, client http.Client, db *sqlx.DB, redis *redis.Client) *Controllers {
	responder := responder.NewResponder(logger)
	return &Controllers{
		SearchController: controller.NewSearcher(logger, client, responder, redis, db),
		AuthController:   controller.NewAuther(logger, client, responder),
		SwagController:   controller.NewSwaggerer(logger, client, responder),
	}
}
