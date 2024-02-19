package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"rpc/service/models"
)

type SearchStorage struct {
	db *sqlx.DB
}

type GeocodeRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

func NewSearchStorage(db *sqlx.DB) *SearchStorage {
	return &SearchStorage{db: db}
}

func (r *SearchStorage) SearchAddress(ctx context.Context, query string) (models.Geo, error) {
	res := models.SearchIntoDb{}
	rows, err := r.db.Queryx("SELECT * FROM address WHERE query=$1 ", query)
	if err != nil {
		return models.Geo{}, err
	}
	for rows.Next() {
		err := rows.StructScan(&res)
		if err != nil {
			continue
		}
		break
	}
	return models.Geo{
		Result: res.Result,
		GeoLat: res.GeoLat,
		GeoLon: res.GeoLon,
	}, nil
}

func (r *SearchStorage) SaveAddress(ctx context.Context, query string, address models.Geo) error {
	_, err := r.db.Exec("INSERT INTO address (query,result,lat,lon) VALUES ($1,$2,$3,$4)", query, address.Result, address.GeoLat, address.GeoLon)
	if err != nil {
		return err
	}
	return nil
}

func (r *SearchStorage) SearchGeoCode(ctx context.Context, request models.GeocodeRequest) (models.Geo, error) {
	res := models.GeoIntoDb{}
	rows, err := r.db.Queryx("SELECT * FROM geo WHERE (lat=$1, lng=$2)", request.Lat, request.Lng)
	if err != nil {
		return models.Geo{}, err
	}
	for rows.Next() {
		err := rows.StructScan(&res)
		if err != nil {
			continue
		}
		break
	}
	return models.Geo{
		Result: res.Result,
		GeoLat: res.GeoLat,
		GeoLon: res.GeoLon,
	}, nil

}

func (r *SearchStorage) SaveGeoCode(ctx context.Context, request models.GeocodeRequest, response models.Geo) error {
	_, err := r.db.Exec("INSERT INTO geo (lat,lng,result,r_lat,r_lon) VALUES ($1,$2,$3,$4,$5)", request.Lat, request.Lng, response.Result, response.GeoLat, response.GeoLon)
	if err != nil {
		return err
	}
	return nil

}
