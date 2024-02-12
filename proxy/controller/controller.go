package controller

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	pb "hugoproxy-main/proxy/proto/gen"
	"hugoproxy-main/proxy/responder"
	service "hugoproxy-main/proxy/service"
	"hugoproxy-main/proxy/storage"
	"net/http"
	"net/rpc"
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

func NewSearcherJsonRpc(responder responder.Responder, rpc *rpc.Client) Searcher {
	return &SearcherJsonRpc{
		responder: responder,
		rpc:       rpc,
	}
}

func NewSearcher(responder responder.Responder, rpc *rpc.Client) Searcher {
	return &Search{
		responder: responder,
		rpc:       rpc,
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

func NewSearchgRpc(responder responder.Responder, cc grpc.ClientConnInterface) *SearchgRpc {
	return &SearchgRpc{
		responder: responder,
		rpc:       pb.NewGeoServiceClient(cc),
	}
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
type SearcherJsonRpc struct {
	responder responder.Responder
	logger    zap.Logger
	rpc       *rpc.Client
}
type Search struct {
	responder responder.Responder
	logger    zap.Logger
	rpc       *rpc.Client
}
type SearchgRpc struct {
	responder responder.Responder
	rpc       pb.GeoServiceClient
}

func (controller *SearchgRpc) GetSearch(w http.ResponseWriter, r *http.Request) {
	var toService RequestGeoSearch
	err := json.NewDecoder(r.Body).Decode(&toService)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
	}
	ctx := context.Background()
	req := &pb.GeoRequest{Query: toService.Query}
	resp, err := controller.rpc.GetGeoResponse(ctx, req)
	if err != nil {
		controller.responder.ErrorInternal(w, fmt.Errorf("Internal error"))
	}
	geoResp := &GeoResponse{Addresses: []Geo{{Result: resp.GetResult(), GeoLat: resp.GetLat(), GeoLon: resp.GetLon()}}}
	json, err := json.Marshal(geoResp)
	controller.responder.OutputJSON(w, string(json))

}
func (controller *SearchgRpc) GetGeo(w http.ResponseWriter, r *http.Request) {
	var toGeo RequestGeoGeo
	err := json.NewDecoder(r.Body).Decode(&toGeo)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
	}
	ctx := context.Background()
	req := &pb.SearchRequest{Lat: toGeo.Lat, Lon: toGeo.Lng}
	resp, err := controller.rpc.GetSearchResponse(ctx, req)
	if err != nil {
		controller.responder.ErrorInternal(w, fmt.Errorf("Internal error"))
	}
	geoResp := &GeoResponse{Addresses: []Geo{{Result: resp.GetResult(), GeoLat: resp.GetLat(), GeoLon: resp.GetLon()}}}
	json, err := json.Marshal(geoResp)
	controller.responder.OutputJSON(w, string(json))

}

type AddresRequest struct {
	Query string `json:"query"`
}
type AddressResopnse struct {
	Addresses []Geo `json:"addresses" db:"addresses"`
}

func (controller *SearcherJsonRpc) GetSearch(w http.ResponseWriter, r *http.Request) {
	var toService RequestGeoSearch
	err := json.NewDecoder(r.Body).Decode(&toService)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
		return
	}
	request := &AddresRequest{Query: toService.Query}
	var response AddressResopnse
	err = controller.rpc.Call("GeoControllerJsonRpc.GetAddress", request, &response)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
	}
	bytes, _ := json.Marshal(response)
	controller.responder.OutputJSON(w, string(bytes))
}

type GeoRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}
type GeoResponse struct {
	Addresses []Geo `json:"addresses"`
}

type Geo struct {
	Result string `json:"result"`
	GeoLat string `json:"lat"`
	GeoLon string `json:"lon" `
}

func (controller *SearcherJsonRpc) GetGeoCode(w http.ResponseWriter, r *http.Request) {
	var toService GeoRequest
	err := json.NewDecoder(r.Body).Decode(&toService)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
		return
	}
	request := &GeoRequest{Lat: toService.Lat, Lng: toService.Lng}
	var result GeoResponse
	err = controller.rpc.Call("GeoControllerJsonRpc.GetGeo", request, &result)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
	}
	bytes, _ := json.Marshal(result)
	controller.responder.OutputJSON(w, string(bytes))
}
func (controller *Search) GetSearch(w http.ResponseWriter, r *http.Request) {
	var toService RequestGeoSearch
	err := json.NewDecoder(r.Body).Decode(&toService)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
		return
	}
	request := toService.Query
	var response string
	err = controller.rpc.Call("GeoController.GetAddress", &request, &response)
	if err != nil {
		if err != nil {
			controller.responder.ErrorInternal(w, err)
		}
	}

	controller.responder.OutputJSON(w, response)

}
func (controller *Search) GetGeoCode(w http.ResponseWriter, r *http.Request) {
	var toService RequestGeoGeo
	err := json.NewDecoder(r.Body).Decode(&toService)
	if err != nil {
		controller.responder.ErrorBadRequest(w, fmt.Errorf("Invalid json"))
		return
	}
	var request []string = []string{toService.Lat, toService.Lng}
	var response string
	err = controller.rpc.Call("GeoController.GetGeo", &request, &response)
	if err != nil {
		controller.responder.ErrorInternal(w, err)
	}
	controller.responder.OutputJSON(w, response)

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
