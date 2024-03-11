package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"proxy/internal/infra/responder"
	pba "proxy/proto/auth"
	pbs "proxy/proto/geo"
	pbu "proxy/proto/user"
	"strings"
)

// ///////////////////////////////////////////////////////AUTH  ENDPOINT HANDLERS//////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type Auth struct {
	responder responder.Responder
	logger    *zap.SugaredLogger
	rpc       pba.AuthServiceClient
}

func NewAuthController(logger *zap.SugaredLogger, responder responder.Responder, cc grpc.ClientConnInterface) *Auth {
	return &Auth{
		responder: responder,
		logger:    logger,
		rpc:       pba.NewAuthServiceClient(cc),
	}
}

func (controller *Auth) Register(w http.ResponseWriter, r *http.Request) {
	var request RequestAuth
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
		return
	}
	ctx := context.Background()
	req := pba.User{Email: request.Email, Phone: request.Phone, Password: request.Password}
	_, err = controller.rpc.Register(ctx, &req)
	if err != nil {
		controller.responder.ErrorInternal(w, err)
		return
	}
	controller.responder.OutputJSON(w, fmt.Sprintf("Sucseful register "))

}

func (controller *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var request RequestAuth
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
		return
	}
	ctx := context.Background()
	resp, err := controller.rpc.Login(ctx, &pba.User{Email: request.Email, Phone: request.Phone, Password: request.Password})
	if err != nil {
		controller.responder.ErrorInternal(w, fmt.Errorf("Invalid username , email  or password"))
	}
	token := fmt.Sprintf("Bearer" + " " + resp.Token)
	controller.responder.OutputJSON(w, token)
}

func (controller *Auth) Verif(token string) bool {
	ctx := context.Background()
	raws := strings.Split(token, " ")
	if len(raws) != 2 {
		return false
	}
	req := pba.Token{
		Token: raws[1],
	}
	isAuth, err := controller.rpc.Authorised(ctx, &req)
	if err != nil {
		return false
	}
	if !isAuth.IsAuthorised {
		return false
	}

	return true

}

// ///////////////////////////////////////////////////////USER ENDPOINT HANDLERS//////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func NewUserController(responder responder.Responder, logger *zap.SugaredLogger, cc grpc.ClientConnInterface) *User {
	return &User{
		responder: responder,
		logger:    logger,
		rpc:       pbu.NewUserServiceClient(cc),
	}
}

type User struct {
	responder responder.Responder
	logger    *zap.SugaredLogger
	rpc       pbu.UserServiceClient
}

func (controller *User) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var profileIN ProfileRequest
	err := json.NewDecoder(r.Body).Decode(&profileIN)
	log.Println(profileIN)
	req := pbu.ProfileRequest{Email: profileIN.Email, Phone: profileIN.Phone}
	user, err := controller.rpc.GetUser(ctx, &req)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid email or phone", err))
		return
	}
	profileOut := ProfileResponse{Email: user.GetEmail(), Phone: user.GetPhone(), Password: user.GetPassword()}
	b, _ := json.Marshal(&profileOut)
	controller.responder.OutputJSON(w, string(b))
}

func (controller *User) List(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	req := pbu.EmptyRequest{}
	users, err := controller.rpc.List(ctx, &req)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Permission denied", err))
		return

	}
	var list ListUser
	for _, user := range users.User {
		list.Users = append(list.Users, UserFromRpc{Email: user.GetEmail(), Phone: user.GetPhone(), Password: user.GetPassword()})
	}
	b, err := json.Marshal(&list)
	controller.responder.OutputJSON(w, string(b))

}

func (contorller *User) SetUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var In SetUserRequest
	json.NewDecoder(r.Body).Decode(&In)
	req := &pbu.User{
		Email:    In.Email,
		Phone:    In.Phone,
		Password: In.Password,
	}
	_, err := contorller.rpc.SetUser(ctx, req)
	if err != nil {
		contorller.responder.ErrorInternal(w, err)
		return
	}
	contorller.responder.OutputJSON(w, fmt.Sprintf("Sucsess add user %s", In.Email))
}

func (controller *User) Profile(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var profileIN ProfileRequest
	err := json.NewDecoder(r.Body).Decode(&profileIN)
	log.Println(profileIN)
	req := pbu.ProfileRequest{Email: profileIN.Email, Phone: profileIN.Phone}
	user, err := controller.rpc.GetUser(ctx, &req)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid email or phone", err))
		return
	}
	profileOut := ProfileResponse{Email: user.GetEmail(), Phone: user.GetPhone(), Password: user.GetPassword()}
	b, _ := json.Marshal(&profileOut)
	controller.responder.OutputJSON(w, string(b))

}

// ///////////////////////////////////////////////////////GEO ENDPOINT HANDLERS//////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func NewSearchController(responder responder.Responder, logger *zap.SugaredLogger, cc grpc.ClientConnInterface) *Search {
	return &Search{
		responder: responder,
		logger:    logger,
		rpc:       pbs.NewGeoServiceClient(cc),
	}
}

type Search struct {
	responder responder.Responder
	logger    *zap.SugaredLogger
	rpc       pbs.GeoServiceClient
}

func (controller *Search) GetSearch(w http.ResponseWriter, r *http.Request) {
	var toService RequestGeoSearch
	err := json.NewDecoder(r.Body).Decode(&toService)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
		return
	}
	ctx := context.Background()
	req := &pbs.AddressRequest{Query: toService.Query}
	resp, err := controller.rpc.GetAddress(ctx, req)
	if err != nil {
		log.Println(err)
		controller.responder.ErrorInternal(w, fmt.Errorf("Internal error"))
		return
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
		return
	}
	ctx := context.Background()
	req := &pbs.GeoRequest{Lat: toGeo.Lat, Lon: toGeo.Lng}
	resp, err := controller.rpc.GetGeo(ctx, req)
	if err != nil {
		log.Println(err)
		controller.responder.ErrorInternal(w, fmt.Errorf("Internal error"))
		return
	}
	geoResp := &GeoResponse{Addresses: []Geo{{Result: resp.GetResult(), GeoLat: resp.GetLat(), GeoLon: resp.GetLon()}}}
	json, err := json.Marshal(geoResp)
	controller.responder.OutputJSON(w, string(json))

}

// ///////////////////////////////////////////////////////SWAGGER ENDPOINT HANDLERS//////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//	func NewSwaggerer(logger *zap.SugaredLogger, client http.Client, responder responder.Responder) *Swagger {
//		return &Swagger{
//			responder:      responder,
//			logger:         logger,
//			swaggerservice: service.NewSwaggerService(),
//		}
//	}
//

//	type Swagger struct {
//		responder      responder.Responder
//		swaggerservice service.SwaggerServiceer
//		logger         *zap.SugaredLogger
//	}

//func (controller *Swagger) GetSwaggerHtml(w http.ResponseWriter, r *http.Request) {
//	html := controller.swaggerservice.GetSwaggerHtml()
//	controller.responder.OutputHtml(w, html)
//}
//func (controller *Swagger) GetSwaggerJson(w http.ResponseWriter, r *http.Request) {
//	json := controller.swaggerservice.GetSwaggerJson()
//	controller.responder.OutputJSON(w, json)
//
//}
