package controller

// Geo search and geocode type
type RequestGeoSearch struct {
	Query string `json:"query"`
}

type RequestGeoGeo struct {
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

// Auth register and login type
type RequestAuth struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// / User endponit types
type ProfileRequest struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type ProfileResponse struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
type ListUser struct {
	Users []UserFromRpc `json:"users"`
}

type UserFromRpc struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
type SetUserRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
