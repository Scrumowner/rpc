package models

type SearchRequest struct {
	Query string `json:"query" db:"query"`
}
type GeocodeRequest struct {
	Lat string `json:"lat" db:"lat"`
	Lng string `json:"lng" db:"lng"`
}
type GeoResponse struct {
	Addresses []Geo `json:"addresses" db:"addresses"`
}
type Geo struct {
	Result string `json:"result" db:"result"`
	GeoLat string `json:"lat" db:"lat"`
	GeoLon string `json:"lon" db:"lon"`
}
