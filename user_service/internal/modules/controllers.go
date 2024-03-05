package modules

import (
	"github.com/jmoiron/sqlx"
	"user/internal/modules/controller"
)

type Controllers struct {
	User *controller.UserController
}

func NewControllers(db *sqlx.DB) *Controllers {
	return &Controllers{
		User: controller.NewUserController(db),
	}
}
