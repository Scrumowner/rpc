package models

type SearchIntoDb struct {
	Query  string `json:"query" db:"query" db_type:"text"`
	Result string `json:"result" db:"result" db_type:"text"`
	GeoLat string `json:"lat" db:"lat" db_type:"text"`
	GeoLon string `json:"lon" db:"lon" db_type:"text"`
}

type GeoIntoDb struct {
	Lat    string `json:"lat" db:"lat" db_type:"text"`
	Lng    string `json:"lng" db:"lng" db_type:"text"`
	Result string `json:"result" db:"result" db_type:"text"`
	GeoLat string `json:"lat" db:"r_lat" db_type:"text"`
	GeoLon string `json:"lon" db:"r_lon" db_type:"text"`
}

type Tabler interface {
	TableName() string
}

func (s *SearchIntoDb) TableName() string {
	return "address"
}

func (g *GeoIntoDb) TableName() string {
	return "geo"
}
