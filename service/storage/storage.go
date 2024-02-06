package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"rpc/service/models"
	"time"
)

type ProxyCache struct {
	cache   *RedisCache
	storage *SearchStorage
}

func NewProxyCache(redis *redis.Client, db *sqlx.DB) *ProxyCache {
	return &ProxyCache{
		cache:   NewRedisCache(redis),
		storage: NewSearchStorage(db),
	}
}

func (p *ProxyCache) SearchAddress(ctx context.Context, query string) (models.Geo, error) {
	result := models.Geo{}
	res, err := p.cache.Get(ctx, query)
	if err == redis.Nil || err != nil {
		result, err = p.storage.SearchAddress(ctx, query)
		if err != nil {
			return models.Geo{}, err
		}
		if err == nil {
			return result, nil
		}
	}
	if err == nil {
		err := json.Unmarshal(res, &result)
		if err != nil {
			return models.Geo{}, err
		}
		return result, nil
	}
	return models.Geo{}, err
}

func (p *ProxyCache) SaveSearchAddress(ctx context.Context, query string, address models.Geo) error {
	err := p.cache.Set(ctx, "query", address, time.Hour)
	if err != nil {
		return err
	}
	err = p.storage.SaveAddress(ctx, query, address)
	if err != nil {

	}
	return nil
}

func (p *ProxyCache) SearchGeoCode(ctx context.Context, request models.GeocodeRequest) (models.Geo, error) {
	result := models.Geo{}
	res, err := p.cache.Get(ctx, fmt.Sprintf("lat=%s lon=%s", request.Lat, request.Lng))
	if err == redis.Nil || err != nil {
		storageRes, storageErr := p.storage.SearchGeoCode(ctx, request)
		if storageErr != nil {
			return models.Geo{}, storageErr
		}
		return storageRes, nil
	}
	if err == nil {
		json.Unmarshal(res, &result)
		return result, nil
	}
	return models.Geo{}, err
}

func (p *ProxyCache) SaveGeoCode(ctx context.Context, request models.GeocodeRequest, response models.Geo) error {
	err := p.cache.Set(ctx, fmt.Sprintf("lat=%s lon=%s", request.Lat, request.Lng), response, time.Hour)
	if err != nil {
		return err
	}
	err = p.storage.SaveGeoCode(ctx, request, response)
	if err != nil {
		return err
	}
	return nil
}
