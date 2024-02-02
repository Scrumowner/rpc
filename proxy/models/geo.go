package models

type Addres struct {
	query  string `json:"query" db:"query"`
	Result string `json:"result" db:"result"`
	GeoLat string `json:"lat" db:"r_lat"`
	GeoLon string `json:"lon" db:"r_lon"`
}

type Geo struct {
	Lat    string `json:"lat" db:"lat"`
	Lng    string `json:"lng" db:"lng"`
	Result string `json:"result" db:"result"`
	GeoLat string `json:"lat" db:"r_lat"`
	GeoLon string `json:"lon" db:"r_lon"`
}
