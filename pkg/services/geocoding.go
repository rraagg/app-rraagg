package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/rraagg/rraagg/config"
)

type (
	// Geocoding provides a client for retrieving geocoding data
	GeocodingClient struct {
		// config stores application configuration
		config *config.Config
		client *http.Client
	}

	GeoCodingResponse struct {
		Result Result `json:"result"`
	}

	Result struct {
		Input          Input          `json:"input"`
		AddressMatches []AddressMatch `json:"addressMatches"`
	}

	Input struct {
		Address   Address   `json:"address"`
		Benchmark Benchmark `json:"benchmark"`
	}

	Benchmark struct {
		IsDefault            bool   `json:"isDefault"`
		BenchmarkDescription string `json:"benchmarkDescription"`
		ID                   string `json:"id"`
		BenchmarkName        string `json:"benchmarkName"`
	}

	Address struct {
		City   string `json:"city"`
		Street string `json:"street"`
		State  string `json:"state"`
	}

	AddressMatch struct {
		TigerLine         TigerLine         `json:"tigerLine"`
		Coordinates       Coordinates       `json:"coordinates"`
		AddressComponents AddressComponents `json:"addressComponents"`
		MatchedAddress    string            `json:"matchedAddress"`
	}

	TigerLine struct {
		Side        string `json:"side"`
		TigerLineId string `json:"tigerLineId"`
	}

	Coordinates struct {
		X float64 `json:"x"` // longitude
		Y float64 `json:"y"` // latitude
	}

	AddressComponents struct {
		Zip             string `json:"zip"`
		Streetname      string `json:"streetName"`
		PreType         string `json:"preType"`
		City            string `json:"city"`
		PreDirection    string `json:"preDirection"`
		SuffixDirection string `json:"suffixDirection"`
		FromAddress     string `json:"fromAddress"`
		State           string `json:"state"`
		SuffixType      string `json:"suffixType"`
		ToAddress       string `json:"toAddress"`
		SuffixQualifier string `json:"suffixQualifier"`
		PreQualifier    string `json:"preQualifier"`
	}
)

// NewGeoCodingClient creates a new GeocodingClient
func NewGeoCodingClient(cfg *config.Config) *GeocodingClient {
	return &GeocodingClient{
		config: cfg,
		client: &http.Client{},
	}
}

func (g *GeocodingClient) GetGeocodingCoordinates(
	ctx echo.Context,
	street string,
	city string,
	state string,
) (*Coordinates, error) {
	ctx.Logger().
		Infof("Getting GeoCode for: street=%s, city=%s, state=%s, url=%s", street, city, state, g.config.Geocode.Hostname)
	v := url.Values{}
	v.Add("street", street)
	v.Add("city", city)
	v.Add("state", state)
	v.Add("benchmark", "2020")
	v.Add("format", "json")
	url := fmt.Sprintf(
		"%s%s",
		g.config.Geocode.Hostname,
		v.Encode(),
	)
	ctx.Logger().Infof("URL: %s", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		ctx.Logger().Errorf("unable to create geocoding request: %s", err)
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	ctx.Logger().Infof("Request: %v", req)

	resp, err := g.client.Do(req)
	ctx.Logger().Infof("Response: %v", resp)

	defer resp.Body.Close()
	var data GeoCodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		ctx.Logger().Errorf("unable to DECODE response: %s", err)
		return nil, err
	}

	ctx.Logger().Infof("Data: %v", data)

	if len(data.Result.AddressMatches) == 0 {
		ctx.Logger().Errorf("no address matches found")
		return nil, fmt.Errorf("no address matches found")
	}
	var addressMatches []AddressMatch
	for _, address := range data.Result.AddressMatches {
		addressMatches = append(addressMatches, address)
		ctx.Logger().Infof("Address: %v", address)
	}
	if len(addressMatches) == 0 {
		ctx.Logger().Errorf("no addresses added to addressMatches")
		return nil, fmt.Errorf("no addresses added to addressMatches")
	}
	if len(addressMatches) > 1 {
		ctx.Logger().Errorf("multiple address matches found")
		return nil, fmt.Errorf("multiple address matches found")
	}
	// TODO Save GeocodeResult to database

	return &addressMatches[0].Coordinates, nil
}
