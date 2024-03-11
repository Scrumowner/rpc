package router

import (
	"github.com/go-chi/chi/v5"
	"proxy/internal/modules"
)

func NewRouter(r *chi.Mux, c *controllers.Controllers) *chi.Mux {
	r.Route("/user", func(r chi.Router) {
		r.Post("/add", c.UserController.SetUser)
		r.Post("/get", c.UserController.GetUser)
		r.Post("/profile", c.UserController.Profile)
		r.Get("/list", c.UserController.List)
	})
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", c.AuthController.Register)
		r.Post("/login", c.AuthController.Login)
	})
	r.Route("/api/address", func(r chi.Router) {
		r.Post("/search", c.SearchController.GetSearch)
		r.Post("/geocode", c.SearchController.GetGeo)
	})

	//r.Route("/swagger", func(r chi.Router) {
	//	r.Get("/index", controller.SwagController.GetSwaggerHtml)
	//	r.Get("/swagger", controller.SwagController.GetSwaggerJson)
	//})
	return r
}
