package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const geoUrl = "http://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address"
const addressUrl = "https://cleaner.dadata.ru/api/v1/clean/address"

const postMethod = "POST"

const tokenAuthorization string = "Token ff49ea7a69ca05be457b7836ecfe9fbb6dac1417"
const tokenXSecret string = "85555da9ba833d12825aced346392b6df0ad5697"

type Providerer interface {
	GetSearchFromApi(query SearchRequest) ReworkedSearchResponse
	GetGeoCodeFromApi(query GeocodeRequest) ReworkedSearchResponse
}

func NewProvider(client http.Client) Providerer {
	return &GeoProvider{client: client}
}

type GeoProvider struct {
	client http.Client
}

type SearchRequest struct {
	Query string `json:"query"`
}
type SearchResponse []Address
type Address struct {
	Result string `json:"result"`
	GeoLat string `json:"geo_lat"`
	GeoLon string `json:"geo_lon"`
}

func (provider *GeoProvider) GetSearchFromApi(query SearchRequest) ReworkedSearchResponse {
	payload := strings.NewReader(fmt.Sprintf(`["%s"]`, query.Query))

	req, err := http.NewRequest(postMethod, addressUrl, payload)
	ctx := context.Background()
	defer ctx.Done()
	if err != nil {
		return ReworkedSearchResponse{}

	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", tokenAuthorization)
	req.Header.Add("X-Secret", tokenXSecret)

	res, err := provider.client.Do(req)
	if err != nil {
		return ReworkedSearchResponse{}
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ReworkedSearchResponse{}
	}
	searchResp := SearchResponse{}
	json.Unmarshal(body, &searchResp)

	return ReworkedSearchResponse{Addresses: []ReworkedSearch{{Result: searchResp[0].Result, GeoLat: searchResp[0].GeoLat, GeoLon: searchResp[0].GeoLon}}}
}

type GeocodeRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

func (provider *GeoProvider) GetGeoCodeFromApi(query GeocodeRequest) ReworkedSearchResponse {

	payload := strings.NewReader(fmt.Sprintf("{ \"lat\":%s , \"lon\":%s }", query.Lat, query.Lng))

	req, err := http.NewRequest(postMethod, geoUrl, payload)
	if err != nil {

		return ReworkedSearchResponse{}
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", tokenAuthorization)
	resp, err := provider.client.Do(req)
	if err != nil {
		return ReworkedSearchResponse{}
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		return ReworkedSearchResponse{}
	}
	geoAddr := GeoAddresses{}
	json.Unmarshal(data, &geoAddr)
	return ReworkedSearchResponse{Addresses: []ReworkedSearch{{Result: geoAddr.Suggestions[0].Value,
		GeoLat: *geoAddr.Suggestions[0].Data["geo_lat"],
		GeoLon: *geoAddr.Suggestions[0].Data["geo_lon"]}}}

}

type GeoAddresses struct {
	Suggestions []Suggestion `json:"suggestions"`
}

type Suggestion struct {
	Value             string             `json:"value"`
	UnrestrictedValue string             `json:"unrestricted_value"`
	Data              map[string]*string `json:"data"`
}

type ReworkedSearchResponse struct {
	Addresses []ReworkedSearch `json:"addresses"`
}
type ReworkedSearch struct {
	Result string `json:"result"`
	GeoLat string `json:"lat"`
	GeoLon string `json:"lon"`
}
