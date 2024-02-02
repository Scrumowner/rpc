package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
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

func (p *ProxyCache) SearchAddress(ctx context.Context, query string) (ReworkedSearch, error) {
	result := ReworkedSearch{}
	res, err := p.cache.Get(ctx, query)
	if err == redis.Nil {
		result, err = p.storage.SearchAddress(ctx, query)
		if err != nil {
			return ReworkedSearch{}, err
		}
	}
	if err != nil {
		result, err = p.storage.SearchAddress(ctx, query)
		if err != nil {
			return ReworkedSearch{}, err
		}
	}
	if err == nil {
		err := json.Unmarshal(res, &result)
		if err != nil {
			return ReworkedSearch{}, err
		}
		return result, nil
	}
	return ReworkedSearch{}, err
}

func (p *ProxyCache) SaveSearchAddress(ctx context.Context, query string, address ReworkedSearch) error {
	err := p.cache.Set(ctx, "query", address, time.Hour)
	if err != nil {
		return err
	}
	err = p.storage.SaveAddress(ctx, query, address)
	if err != nil {

	}
	return nil
}

func (p *ProxyCache) SearchGeoCode(ctx context.Context, request GeocodeRequest) (ReworkedSearch, error) {
	result := ReworkedSearch{}
	res, err := p.cache.Get(ctx, fmt.Sprintf("lat=%s lon=%s", request.Lat, request.Lng))
	if err == redis.Nil {
		storageRes, storageErr := p.storage.SearchGeoCode(ctx, request)
		if storageErr != nil {
			return ReworkedSearch{}, storageErr
		}
		return storageRes, nil
	}
	if err != nil {
		storageRes, storageErr := p.storage.SearchGeoCode(ctx, request)
		if storageErr != nil {
			return ReworkedSearch{}, err
		}
		return storageRes, nil
	}
	if err == nil {
		json.Unmarshal(res, &result)
		return result, nil
	}
	return ReworkedSearch{}, err
}

func (p *ProxyCache) SaveGeoCode(ctx context.Context, request GeocodeRequest, response ReworkedSearch) error {
	err := p.cache.Set(ctx, fmt.Sprintf("lat=%s lon=%s", request.Lat, request.Lng), response, time.Hour)
	if err != nil {
		return err
	}
	err = p.SaveGeoCode(ctx, request, response)
	if err != nil {
		return err
	}
	return nil
}
