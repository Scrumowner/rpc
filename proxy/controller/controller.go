package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	pba "hugoproxy-main/proxy/proto/auth"
	pb "hugoproxy-main/proxy/proto/geo"
	pbu "hugoproxy-main/proxy/proto/user"
	"hugoproxy-main/proxy/responder"
	service "hugoproxy-main/proxy/service"
	"net/http"
)

func NewSwaggerer(logger *zap.SugaredLogger, client http.Client, responder responder.Responder) *Swagger {
	return &Swagger{
		responder:      responder,
		logger:         logger,
		swaggerservice: service.NewSwaggerService(),
	}
}

func NewAuthController(logger *zap.SugaredLogger, client http.Client, responder responder.Responder, cc grpc.ClientConnInterface) *Auth {
	return &Auth{
		responder: responder,
		logger:    logger,
		rpc:       pba.NewAuthServiceClient(cc),
	}
}

func NewSearchController(responder responder.Responder, logger *zap.SugaredLogger, cc grpc.ClientConnInterface) *Search {
	return &Search{
		responder: responder,
		logger:    logger,
		rpc:       pb.NewGeoServiceClient(cc),
	}
}
func NewUserController(responder responder.Responder, logger *zap.SugaredLogger, cc grpc.ClientConnInterface) *User {
	return &User{
		responder: responder,
		logger:    logger,
		rpc:       pbu.NewUserServiceClient(cc),
	}
}

type Swagger struct {
	responder      responder.Responder
	swaggerservice service.SwaggerServiceer
	logger         *zap.SugaredLogger
}
type Auth struct {
	responder responder.Responder
	logger    *zap.SugaredLogger
	rpc       pba.AuthServiceClient
}
type User struct {
	responder responder.Responder
	logger    *zap.SugaredLogger
	rpc       pbu.UserServiceClient
}
type Search struct {
	responder responder.Responder
	logger    *zap.SugaredLogger
	rpc       pb.GeoServiceClient
}

func (controller *Search) GetSearch(w http.ResponseWriter, r *http.Request) {
	var toService RequestGeoSearch
	err := json.NewDecoder(r.Body).Decode(&toService)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
	}
	ctx := context.Background()
	req := &pb.AddressRequest{Query: toService.Query}
	resp, err := controller.rpc.GetAddress(ctx, req)
	if err != nil {
		controller.responder.ErrorInternal(w, fmt.Errorf("Internal error"))
	}
	geoResp := &GeoResponse{Addresses: []Geo{{Result: resp.GetResult(), GeoLat: resp.GetLat(), GeoLon: resp.GetLon()}}}
	json, err := json.Marshal(geoResp)
	controller.responder.OutputJSON(w, string(json))

}
func (controller *Search) GetGeo(w http.ResponseWriter, r *http.Request) {
	var toGeo RequestGeoGeo
	err := json.NewDecoder(r.Body).Decode(&toGeo)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
	}
	ctx := context.Background()
	req := &pb.GeoRequest{Lat: toGeo.Lat, Lon: toGeo.Lng}
	resp, err := controller.rpc.GetGeo(ctx, req)
	if err != nil {
		controller.responder.ErrorInternal(w, fmt.Errorf("Internal error"))
	}
	geoResp := &GeoResponse{Addresses: []Geo{{Result: resp.GetResult(), GeoLat: resp.GetLat(), GeoLon: resp.GetLon()}}}
	json, err := json.Marshal(geoResp)
	controller.responder.OutputJSON(w, string(json))

}

func (controller *Auth) Register(w http.ResponseWriter, r *http.Request) {
	var request RequestAuth
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
		return
	}
	ctx := context.Background()
	resp, err := controller.rpc.Register(ctx, &pba.User{Email: request.Email, Password: request.Password})
	if err != nil {
		controller.responder.ErrorInternal(w, err)
	}
	controller.responder.OutputJSON(w, resp.Response)

}

func (controller *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var request RequestAuth
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
		return
	}
	ctx := context.Background()
	resp, err := controller.rpc.Login(ctx, &pba.User{Email: request.Email, Password: request.Password})
	if err != nil {
		controller.responder.ErrorInternal(w, fmt.Errorf("Invalid username or password"))
	}
	if err != nil {
		controller.responder.ErrorUnauthorized(w, err)
	}
	controller.responder.OutputJSON(w, resp.Token)
}
func (controller *Auth) Verif(token string) bool {
	ctx := context.Background()
	isAuth, err := controller.rpc.Authorised(ctx, &pba.Token{Token: token})
	if err != nil {
		return false
	}
	if !isAuth.IsAuthorised {
		return false
	}

	return true

}

func (controller *User) Profile(w http.ResponseWriter, r *http.Request) {
	var profileIN ProfileRequest
	err := json.NewDecoder(r.Body).Decode(&profileIN)
	ctx := context.Background()
	user, err := controller.rpc.Profile(ctx, &pbu.ProfileRequest{Email: profileIN.Email})
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid email"))
	}
	var profileOut ProfileResponse = ProfileResponse{Email: user.GetEmail(), Password: user.GetPassword()}
	b, _ := json.Marshal(&profileOut)
	controller.responder.OutputJSON(w, string(b))

}

func (controller *User) List(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	users, err := controller.rpc.List(ctx, &pbu.EmptyRequest{})
	if err != nil {
		if err != nil {
			controller.responder.ErrorBadRequest(w, fmt.Errorf("Permission denied"))
		}
	}
	var list LisetResponse
	for _, user := range users.User {
		list.Users = append(list.Users, ListUser{Email: user.GetEmail(), Password: user.GetPassword()})
	}
	b, err := json.Marshal(&list)
	controller.responder.OutputJSON(w, string(b))

}
func (controller *Swagger) GetSwaggerHtml(w http.ResponseWriter, r *http.Request) {
	html := controller.swaggerservice.GetSwaggerHtml()
	controller.responder.OutputHtml(w, html)
}
func (controller *Swagger) GetSwaggerJson(w http.ResponseWriter, r *http.Request) {
	json := controller.swaggerservice.GetSwaggerJson()
	controller.responder.OutputJSON(w, json)

}
