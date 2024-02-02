package gen

import (
	"hugoproxy-main/proxy/provider"
	"hugoproxy-main/proxy/storage"
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

// swagger:route POST /api/address/search apiRequestSearch
// Search by address
// responses:
//
//	200: apiResponseSearch
//
// swagger:parameters apiRequestSearch
type apiRequestSearch struct {
	// in:body
	// required: true
	Body provider.SearchRequest

	// in: header
	// name: Authorization
	// description: Bearer token for authentication
	Token string `json:"Authorization"`
}

// swagger:response apiResponseSearch
type apiResponseSearch struct {
	// in:body
	// required: true
	Body provider.ReworkedSearch
}

// swagger:route POST /api/address/geocode apiGeocodeRequest
// Search by geocode
// responses:
//
//	200: apiGeocodeResponse
//
// swagger:parameters apiGeocodeRequest
type apiGeocodeRequest struct {
	// in:body
	// required: true
	Body provider.GeocodeRequest

	// in: header
	// name: Authorization
	// description: Bearer token for authentication
	Token string `json:"Authorization"`
}

// swagger:response apiGeocodeResponse
type apiGeocodeResponse struct {
	// in:body
	// required: true
	Body provider.ReworkedSearchResponse
}

// swagger:route POST /api/register apiRegisterRequest
// Register
// responses:
//
//	200: apiRegisterResponse
//
// swagger:parameters apiRegisterRequest
type apiRegisterRequest struct {
	// in:body
	// required: true
	Body storage.User

	// in: header
	// name: Authorization
	// description: Bearer token for authentication
	Token string `json:"Authorization"`
}

// swagger:response apiRegisterResponse
type apiRegisterResponse struct {
	// in:body
	// required: true
	Body struct {
		Message string `json:"message"`
	}
}

// swagger:route POST /api/login apiLoginRequest
// Login
// responses:
//
//	200: apiLoginResponse
//
// swagger:parameters apiLoginRequest
type apiLoginRequest struct {
	// in:body
	// required: true
	Body storage.User
}

// swagger:response apiLoginResponse
type apiLoginResponse struct {
	// in:body
	// required: true
	Body struct {
		Token string `json:"token"`
	}
}

// swagger:securityDefinitions
//   Bearer:
//     type: apiKey
//     name: Authorization
//     in: header
