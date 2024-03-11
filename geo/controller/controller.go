package controller

import (
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"rpc/service/config"
	"rpc/service/models"
	pb "rpc/service/proto/geo"
	"rpc/service/service"
)

type GeoControllergRpc struct {
	Serivce service.GeoServiceer
	pb.UnimplementedGeoServiceServer
}

func NewGeoContollergRpc(client http.Client, redis *redis.Client, db *sqlx.DB, cfg *config.Config) *GeoControllergRpc {
	return &GeoControllergRpc{
		Serivce: service.NewGeoService(client, redis, db, cfg),
	}
}

func (s *GeoControllergRpc) GetAddress(ctx context.Context, req *pb.AddressRequest) (*pb.GeoResponse, error) {
	geo, err := s.Serivce.GetSearch(models.SearchRequest{Query: req.GetQuery()})
	if err != nil {
		log.Println(err)
		return &pb.GeoResponse{}, err
	}
	resp := models.GeoResponse{}

	err = json.Unmarshal([]byte(geo), &resp)
	if err != nil {
		log.Println(err)
		return &pb.GeoResponse{}, err
	}
	return &pb.GeoResponse{
		Result: resp.Addresses[0].Result,
		Lat:    resp.Addresses[0].GeoLat,
		Lon:    resp.Addresses[0].GeoLon,
	}, nil
}

func (s *GeoControllergRpc) GetGeo(ctx context.Context, req *pb.GeoRequest) (*pb.GeoResponse, error) {
	geo, err := s.Serivce.GetGeoCode(models.GeocodeRequest{Lat: req.GetLat(), Lng: req.GetLon()})
	if err != nil {
		log.Println(err)
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

// ///////////////////////////////////JSONRPC GEO CONTROLLER //////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//// /types for json rpc calls
//type AddresRequest struct {
//	Query string `json:"query"`
//}
//type AddressResopnse struct {
//	Addresses []Geo `json:"addresses" db:"addresses"`
//}
//
//// /type for json rpc calls
//type GeoRequest struct {
//	Lat string `json:"lat"`
//	Lng string `json:"lng"`
//}
//type GeoResponse struct {
//	Addresses []Geo `json:"addresses"`
//}
//
//type Geo struct {
//	Result string `json:"result"`
//	GeoLat string `json:"lat"`
//	GeoLon string `json:"lon" `
//}
//
//type GeoControllerJsonRpc struct {
//	service service.GeoServiceer
//}
//
//func NewGeoControllerJsonRpc(client http.Client, redis *redis.Client, db *sqlx.DB, cfg *config.Config) *GeoControllerJsonRpc {
//	return &GeoControllerJsonRpc{
//		service: service.NewGeoService(client, redis, db, cfg),
//	}
//
//}
//
//func NewGeoController(client http.Client, redis *redis.Client, db *sqlx.DB, cfg *config.Config) *GeoController {
//	return &GeoController{
//		service: service.NewGeoService(client, redis, db, cfg),
//	}
//
//}
//
//func (g *GeoControllerJsonRpc) GetAddress(arg *AddresRequest, resp *AddressResopnse) error {
//	res, err := g.service.GetSearch(models.SearchRequest{Query: arg.Query})
//	if err != nil {
//		*resp = AddressResopnse{}
//		return nil
//	}
//	json.Unmarshal([]byte(res), resp)
//	return nil
//}
//func (g *GeoControllerJsonRpc) GetGeo(arg *GeoRequest, resp *GeoResponse) error {
//	res, err := g.service.GetGeoCode(models.GeocodeRequest{Lat: arg.Lat, Lng: arg.Lng})
//	if err != nil {
//		*resp = GeoResponse{}
//		return nil
//	}
//	json.Unmarshal([]byte(res), resp)
//	return nil
//}
//
//////////////////////////////////////////////RPC GEO CONTROLLER /////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//type GeoControllerer interface {
//	GetAddress(args *string, reply *string) error
//	GetGeo(args []string, reply *string) error
//}
//
//type GeoController struct {
//	service service.GeoServiceer
//}
//
//func (g *GeoController) GetAddress(args *string, reply *string) error {
//	res, err := g.service.GetSearch(models.SearchRequest{Query: *args})
//	if err != nil {
//		*reply = fmt.Sprintf("Error")
//	}
//	*reply = res
//	return nil
//}
//
//func (g *GeoController) GetGeo(args []string, reply *string) error {
//	lat, lng := args[0], args[1]
//	res, err := g.service.GetGeoCode(models.GeocodeRequest{Lat: lat, Lng: lng})
//	if err != nil {
//		*reply = "Erorr"
//	}
//	*reply = res
//	return nil
//
//}
