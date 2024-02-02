package controller

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"hugoproxy-main/proxy/provider"
	"hugoproxy-main/proxy/responder"
	service "hugoproxy-main/proxy/service"
	"hugoproxy-main/proxy/storage"
	"net/http"
)

type Searcher interface {
	GetSearch(w http.ResponseWriter, r *http.Request)
	GetGeoCode(w http.ResponseWriter, r *http.Request)
}
type Auther interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}
type Swaggerer interface {
	GetSwaggerHtml(w http.ResponseWriter, r *http.Request)
	GetSwaggerJson(w http.ResponseWriter, r *http.Request)
}

func NewSearcher(logger *zap.SugaredLogger, client http.Client, responder responder.Responder, redis *redis.Client, db *sqlx.DB) Searcher {
	return &Search{
		responder:  responder,
		geoservice: service.NewGeoService(provider.NewProvider(client), redis, db),
	}
}
func NewAuther(logger *zap.SugaredLogger, client http.Client, responder responder.Responder) Auther {
	return &Auth{
		responder:   responder,
		authservice: service.NewAuthService(storage.NewStorage()),
	}
}
func NewSwaggerer(logger *zap.SugaredLogger, client http.Client, responder responder.Responder) Swaggerer {
	return &Swagger{
		responder:      responder,
		swaggerservice: service.NewSwaggerService(),
	}
}

type Search struct {
	responder  responder.Responder
	geoservice service.GeoServiceer
	logger     zap.Logger
}
type Auth struct {
	responder   responder.Responder
	authservice service.AuthServiceer
	logger      zap.Logger
}
type Swagger struct {
	responder      responder.Responder
	swaggerservice service.SwaggerServiceer
	logger         zap.Logger
}

func (controller *Search) GetSearch(w http.ResponseWriter, r *http.Request) {
	var toService RequestGeoSearch
	err := json.NewDecoder(r.Body).Decode(&toService)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
		return
	}
	json, err := controller.geoservice.GetSearch(service.SearchRequestService{Query: toService.Query})
	if err != nil {
		controller.responder.ErrorInternal(w, err)
	}
	controller.responder.OutputJSON(w, json)

}
func (controller *Search) GetGeoCode(w http.ResponseWriter, r *http.Request) {
	var toService RequestGeoGeo
	err := json.NewDecoder(r.Body).Decode(&toService)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
		return
	}
	json, err := controller.geoservice.GetGeoCode(service.GeocodeRequestService{Lat: toService.Lat, Lng: toService.Lng})
	if err != nil {
		controller.responder.ErrorInternal(w, err)
	}
	controller.responder.OutputJSON(w, json)

}
func (controller *Auth) Register(w http.ResponseWriter, r *http.Request) {
	var request RequestAuth
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
		return
	}
	output, err := controller.authservice.Register(service.UserService{ID: request.ID,
		Name:     request.Name,
		Phone:    request.Phone,
		Email:    request.Email,
		Password: request.Password})
	if err != nil {
		controller.responder.ErrorInternal(w, err)
	}
	controller.responder.OutputJSON(w, output)

}
func (controller *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var request RequestAuth
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
		return
	}
	token, err := controller.authservice.Login(service.UserService{ID: request.ID,
		Name:     request.Name,
		Phone:    request.Phone,
		Email:    request.Email,
		Password: request.Password})
	if err != nil {
		controller.responder.ErrorUnauthorized(w, err)
	}
	controller.responder.OutputJSON(w, token)
}
func (controller *Swagger) GetSwaggerHtml(w http.ResponseWriter, r *http.Request) {
	html := controller.swaggerservice.GetSwaggerHtml()
	controller.responder.OutputHtml(w, html)
}
func (controller *Swagger) GetSwaggerJson(w http.ResponseWriter, r *http.Request) {
	json := controller.swaggerservice.GetSwaggerJson()
	controller.responder.OutputJSON(w, json)

}

type RequestGeoSearch struct {
	Query string `json:"query"`
}
type RequestGeoGeo struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type RequestAuth struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
