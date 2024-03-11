package service

import (
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"rpc/service/config"
	"rpc/service/models"
	"rpc/service/provider"
	"rpc/service/storage"
	"strings"
)

type GeoServiceer interface {
	GetSearch(r models.SearchRequest) (string, error)
	GetGeoCode(r models.GeocodeRequest) (string, error)
}

type GeoService struct {
	provider     provider.Providerer
	storageProxy *storage.ProxyCache
}

func NewGeoService(client http.Client, redis *redis.Client, db *sqlx.DB, cfg *config.Config) GeoServiceer {
	return &GeoService{provider: provider.NewProvider(client, cfg), storageProxy: storage.NewProxyCache(redis, db)}
}

func (geo *GeoService) GetSearch(r models.SearchRequest) (string, error) {
	var result = strings.Builder{}
	res, err := geo.storageProxy.SearchAddress(context.Background(), r.Query)
	if err == nil {
		json.NewEncoder(&result).Encode(models.GeoResponse{Addresses: []models.Geo{res}})

		return result.String(), nil
	}

	resp := geo.provider.GetSearchFromApi(provider.SearchRequest{Query: r.Query})
	err = json.NewEncoder(&result).Encode(models.GeoResponse{Addresses: []models.Geo{{
		Result: resp.Addresses[0].Result,
		GeoLat: resp.Addresses[0].GeoLat,
		GeoLon: resp.Addresses[0].GeoLon,
	}}})
	if err != nil {
		return "", err
		log.Println(err)
	}
	geo.storageProxy.SaveSearchAddress(context.Background(), r.Query, models.Geo{
		Result: resp.Addresses[0].Result,
		GeoLat: resp.Addresses[0].GeoLat,
		GeoLon: resp.Addresses[0].GeoLon,
	})
	return result.String(), nil

}
func (geo *GeoService) GetGeoCode(r models.GeocodeRequest) (string, error) {
	var result = strings.Builder{}
	res, err := geo.storageProxy.SearchGeoCode(context.Background(), models.GeocodeRequest{Lat: r.Lat, Lng: r.Lng})
	if err == nil {
		json.NewEncoder(&result).Encode(models.GeoResponse{Addresses: []models.Geo{res}})
		return result.String(), nil
	}
	respFromProvider := geo.provider.GetGeoCodeFromApi(provider.GeocodeRequest{Lat: r.Lat, Lng: r.Lng})
	err = json.NewEncoder(&result).Encode(models.GeoResponse{Addresses: []models.Geo{models.Geo{
		Result: respFromProvider.Addresses[0].Result,
		GeoLat: respFromProvider.Addresses[0].GeoLat,
		GeoLon: respFromProvider.Addresses[0].GeoLon}}})
	if err != nil {
		log.Println(err)
		return "", err
	}
	geo.storageProxy.SaveGeoCode(context.Background(), r, models.Geo{
		Result: respFromProvider.Addresses[0].Result,
		GeoLat: respFromProvider.Addresses[0].GeoLat,
		GeoLon: respFromProvider.Addresses[0].GeoLon})
	return result.String(), nil

}
