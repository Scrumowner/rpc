package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"net/http"
	"rpc/service/models"
	pb "rpc/service/proto/gen"
	"rpc/service/service"
)

// /types for json rpc calls
type AddresRequest struct {
	Query string `json:"query"`
}
type AddressResopnse struct {
	Addresses []Geo `json:"addresses" db:"addresses"`
}

// /type for json rpc calls
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

//	type GeoControllererJsonRpc interface {
//		GetAddress(arg *AddresRequest, resp *AddressResopnse) error
//		GetGeo(arg *GeoRequest, resp *GeoResponse) error
//	}
type GeoControllerer interface {
	GetAddress(args *string, reply *string) error
	GetGeo(args []string, reply *string) error
}
type GeoController struct {
	service service.GeoServiceer
}

type GeoControllerJsonRpc struct {
	service service.GeoServiceer
}
type GeoControllergRpc struct {
	Serivce service.GeoServiceer
	pb.UnimplementedGeoServiceServer
}

func NewGeoContollergRpc(client http.Client, redis *redis.Client, db *sqlx.DB) *GeoControllergRpc {
	return &GeoControllergRpc{
		Serivce: service.NewGeoService(client, redis, db),
	}
}
func (s *GeoControllergRpc) GetSearchResponse(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	geo, err := s.Serivce.GetGeoCode(models.GeocodeRequest{Lat: req.GetLat(), Lng: req.GetLon()})
	if err != nil {
		return &pb.SearchResponse{}, err
	}
	resp := models.GeoResponse{}
	json.Unmarshal([]byte(geo), &resp)
	return &pb.SearchResponse{
		Result: resp.Addresses[0].Result,
		Lat:    resp.Addresses[0].GeoLat,
		Lon:    resp.Addresses[0].GeoLon,
	}, nil
}

func (s *GeoControllergRpc) GetGeoResponse(ctx context.Context, req *pb.GeoRequest) (*pb.GeoResponse, error) {
	geo, err := s.Serivce.GetSearch(models.SearchRequest{Query: req.GetQuery()})
	if err != nil {
		return &pb.GeoResponse{}, err
	}
	resp := models.GeoResponse{}
	json.Unmarshal([]byte(geo), &resp)
	return &pb.GeoResponse{
		Result: resp.Addresses[0].Result,
		Lat:    resp.Addresses[0].GeoLat,
		Lon:    resp.Addresses[0].GeoLon,
	}, nil
}
func (s *GeoControllergRpc) mustEmbedUnimplementedGeoServiceServer() {

}

func NewGeoControllerJsonRpc(client http.Client, redis *redis.Client, db *sqlx.DB) *GeoControllerJsonRpc {
	return &GeoControllerJsonRpc{
		service: service.NewGeoService(client, redis, db),
	}

}
func NewGeoController(client http.Client, redis *redis.Client, db *sqlx.DB) *GeoController {
	return &GeoController{
		service: service.NewGeoService(client, redis, db),
	}

}

func (g *GeoControllerJsonRpc) GetAddress(arg *AddresRequest, resp *AddressResopnse) error {
	res, err := g.service.GetSearch(models.SearchRequest{Query: arg.Query})
	if err != nil {
		*resp = AddressResopnse{}
		return nil
	}
	json.Unmarshal([]byte(res), resp)
	return nil
}
func (g *GeoControllerJsonRpc) GetGeo(arg *GeoRequest, resp *GeoResponse) error {
	res, err := g.service.GetGeoCode(models.GeocodeRequest{Lat: arg.Lat, Lng: arg.Lng})
	if err != nil {
		*resp = GeoResponse{}
		return nil
	}
	json.Unmarshal([]byte(res), resp)
	return nil
}

func (g *GeoController) GetAddress(args *string, reply *string) error {
	res, err := g.service.GetSearch(models.SearchRequest{Query: *args})
	if err != nil {
		*reply = fmt.Sprintf("Error")
	}
	*reply = res
	return nil
}
func (g *GeoController) GetGeo(args []string, reply *string) error {
	lat, lng := args[0], args[1]
	res, err := g.service.GetGeoCode(models.GeocodeRequest{Lat: lat, Lng: lng})
	if err != nil {
		*reply = "Erorr"
	}
	*reply = res
	return nil

}
