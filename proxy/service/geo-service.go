package service

import (
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"hugoproxy-main/proxy/provider"
	"hugoproxy-main/proxy/storage"
	"strings"
)

type GeoServiceer interface {
	GetSearch(r SearchRequestService) (string, error)
	GetGeoCode(r GeocodeRequestService) (string, error)
}

type GeoService struct {
	provider     provider.Providerer
	storageProxy *storage.ProxyCache
}

func NewGeoService(provider provider.Providerer, redis *redis.Client, db *sqlx.DB) GeoServiceer {
	return &GeoService{provider: provider, storageProxy: storage.NewProxyCache(redis, db)}
}

func (geo *GeoService) GetSearch(r SearchRequestService) (string, error) {
	var result = strings.Builder{}
	res, err := geo.storageProxy.SearchAddress(context.Background(), r.Query)
	if err == nil {
		json.NewEncoder(&result).Encode(ReworkedSearchResponseService{Addresses: []ReworkedSearchService{{
			Result: res.Result,
			GeoLat: res.GeoLat,
			GeoLon: res.GeoLon,
		}}})

		return result.String(), nil
	}

	resp := geo.provider.GetSearchFromApi(provider.SearchRequest{Query: r.Query})
	err = json.NewEncoder(&result).Encode(ReworkedSearchResponseService{Addresses: []ReworkedSearchService{{
		Result: resp.Addresses[0].Result,
		GeoLat: resp.Addresses[0].GeoLat,
		GeoLon: resp.Addresses[0].GeoLon,
	}}})
	if err != nil {
		return "", err
	}
	geo.storageProxy.SaveSearchAddress(context.Background(), r.Query, storage.ReworkedSearch{
		Result: resp.Addresses[0].Result,
		GeoLat: resp.Addresses[0].GeoLat,
		GeoLon: resp.Addresses[0].GeoLon,
	})
	return result.String(), nil

}
func (geo *GeoService) GetGeoCode(r GeocodeRequestService) (string, error) {
	var result = strings.Builder{}
	res, err := geo.storageProxy.SearchGeoCode(context.Background(), storage.GeocodeRequest{Lat: r.Lat, Lng: r.Lng})
	if err == nil {
		json.NewEncoder(&result).Encode(ReworkedSearchResponseService{Addresses: []ReworkedSearchService{{
			Result: res.Result,
			GeoLat: res.GeoLat,
			GeoLon: res.GeoLon}}})
		return result.String(), nil
	}
	if err != nil {
		respFromProvider := geo.provider.GetGeoCodeFromApi(provider.GeocodeRequest{Lat: r.Lat, Lng: r.Lng})
		err = json.NewEncoder(&result).Encode(ReworkedSearchResponseService{Addresses: []ReworkedSearchService{{
			Result: respFromProvider.Addresses[0].Result,
			GeoLat: respFromProvider.Addresses[0].GeoLat,
			GeoLon: respFromProvider.Addresses[0].GeoLon}}})
		if err != nil {
			return "", err
		}
		geo.storageProxy.SaveGeoCode(context.Background(), storage.GeocodeRequest{Lat: r.Lat, Lng: r.Lng}, storage.ReworkedSearch{Result: respFromProvider.Addresses[0].Result,
			GeoLat: respFromProvider.Addresses[0].GeoLat,
			GeoLon: respFromProvider.Addresses[0].GeoLon})
		return result.String(), nil
	}
	return "", err
}

type SearchRequestService struct {
	Query string `json:"query"`
}
type GeocodeRequestService struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}
type ReworkedSearchResponseService struct {
	Addresses []ReworkedSearchService `json:"addresses"`
}
type ReworkedSearchService struct {
	Result string `json:"result"`
	GeoLat string `json:"lat"`
	GeoLon string `json:"lon"`
}
