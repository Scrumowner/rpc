package gen

import (
	"proxy/internal/modules/controller"
)

// Package classification infoblog.
//
// Documentation of your project API.
//
// Schemes: http, https
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
//
// Consumes:
//   - application/json
//
// Produces:
//   - application/json
//
// Security:
//   - basic
//   - Bearer: []
//
// SecurityDefinitions:
//
//	Bearer:
//	  type: apiKey
//	  name: Authorization
//	  in: header
//
// swagger:meta
//
//go:generate swagger generate spec -o ./swagger.json --scan-models

// swagger:route POST /api/address/search Search
// Search by address
// responses:
//
//	200: apiResponseSearch
//
// swagger:parameters apiRequestSearch
type apiRequestSearch struct {
	// in:body
	// required: true
	Body controller.RequestGeoSearch

	// in: header
	// name: Authorization
	// description: Bearer token for authentication
	Token string `json:"Authorization"`
}

// swagger:response apiResponseSearch
type apiResponseSearch struct {
	// in:body
	// required: true
	Body controller.GeoResponse
}

// swagger:route POST /api/address/geocode Geocode
// Search by geocode
// responses:
//
//	200: apiGeocodeResponse
//
// swagger:parameters apiGeocodeRequest
type apiGeocodeRequest struct {
	// in:body
	// required: true
	Body controller.RequestGeoGeo

	// in: header
	// name: Authorization
	// description: Bearer token for authentication
	Token string `json:"Authorization"`
}

// swagger:response apiGeocodeResponse
type apiGeocodeResponse struct {
	// in:body
	// required: true
	Body controller.GeoResponse
}

// swagger:route POST /auth/register Register
// Register
// responses:
//
//	200: apiRegisterResponse
//
// swagger:parameters apiRegisterRequest
type apiRegisterRequest struct {
	// in:body
	// required: true
	Body controller.RequestAuth
}

// swagger:response apiRegisterResponse
type apiRegisterResponse struct {
	// in:body
	// required: true
	Body struct {
		Message string `json:"message"`
	}
}

// swagger:route POST /auth/login Login
// Login
// responses:
//
//	200: apiLoginResponse
//
// swagger:parameters apiLoginRequest
type apiLoginRequest struct {
	// in:body
	// required: true
	Body controller.RequestAuth
}

// swagger:response apiLoginResponse
type apiLoginResponse struct {
	// in:body
	// required: true
	Body struct {
		Token string `json:"token"`
	}
}

// swagger:route POST /user/profile Profile
// Profile
// responses:
//
// 200 : apiProfileResponse
//
// swagger:parameters apiProfileRequest
type apiProfileRequest struct {
	//in:body
	//required: true
	Body controller.ProfileRequest
	// in: header
	// name: Authorization
	// description: Bearer token for authentication
	Token string `json:"Authorization"`
}

// swagger:response apiProfileResponse
type apiProfileResponse struct {
	//in:body
	//required:true
	Body controller.ProfileResponse
}

// swagger:route POST /user/list List
// List
// responses:
//
// 200 : apiListResponse
//
// swagger:parameters apiListRequest
type apiListRequest struct {
	// in: header
	// name: Authorization
	// description: Bearer token for authentication
	Token string `json:"Authorization"`
}

// swagger:response apiListResponse
type apiListResponse struct {
	//in:body
	//required:true
	Body controller.ListUser
}
