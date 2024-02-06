package controller

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"net/http"
	"rpc/service/models"
	"rpc/service/service"
)

type GeoControllerer interface {
	GetAddress(args *string, reply *string) error
	GetGeo(args []string, reply *string) error
}
type GeoController struct {
	service service.GeoServiceer
}

func NewGeoController(client http.Client, redis *redis.Client, db *sqlx.DB) GeoControllerer {
	return &GeoController{
		service: service.NewGeoService(client, redis, db),
	}

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
