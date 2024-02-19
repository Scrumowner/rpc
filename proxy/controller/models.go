package controller

type RequestGeoSearch struct {
	Query string `json:"query"`
}

type RequestGeoGeo struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}
type RequestAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AddresRequest struct {
	Query string `json:"query"`
}
type AddressResopnse struct {
	Addresses []Geo `json:"addresses"`
}

type GeoRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}
type GeoResponse struct {
	Addresses []Geo `json:"addresses"`
}

type Geo struct {
	Result string `json:"result"`
	GeoLat string `json:"lat"`
	GeoLon string `json:"lon" `
}
type ProfileRequest struct {
	Email string `json:"email"`
}
type ProfileResponse struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LisetResponse struct {
	Users []ListUser `json:"users"`
}
type ListUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
