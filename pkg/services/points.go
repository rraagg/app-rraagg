package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rraagg/rraagg/config"
)

type (
	PointsClient struct {
		config *config.Config
		client *http.Client
	}

	PointsResponse struct {
		Properties PointsProperties `json:"properties"`
	}

	PointsProperties struct {
		ID             string `json:"@id"`
		Type           string `json:"@type"`
		CWA            string `json:"cwa"`
		ForecastOffice string `json:"forecastOffice"`
		GridID         string `json:"gridId"`
		GridX          int    `json:"gridX"`
		GridY          int    `json:"gridY"`
		Forecast       string `json:"forecast"`
		ForecastHourly string `json:"forecastHourly"`
	}
)

// NewPointsClient creates a new PointsClient
func NewPointsClient(cfg *config.Config) *PointsClient {
	return &PointsClient{
		config: cfg,
		client: &http.Client{},
	}
}

func (p *PointsClient) GetPoints(
	ctx echo.Context,
	latitude float64,
	longitude float64,
) (PointsProperties, error) {
	url := fmt.Sprintf("%s%f,%f", p.config.Points.Hostname, latitude, longitude)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		ctx.Logger().Errorf("unable to create points request: %s", err)
		return PointsProperties{}, err
	}
	req.Header.Set("Accept", "application/json")

	ctx.Logger().Infof("Request: %v", req)

	resp, err := p.client.Do(req)
	ctx.Logger().Infof("Response: %v", resp)

	defer resp.Body.Close()
	var data PointsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		ctx.Logger().Errorf("unable to DECODE points response: %s", err)
		return PointsProperties{}, err
	}

	ctx.Logger().Infof("Data: %v", data)
	if err != nil {
	}
	return data.Properties, nil
}
