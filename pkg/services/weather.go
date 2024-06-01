package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rraagg/rraagg/config"
)

type (
	// WeatherClient provides a client for retrieving weather data
	WeatherClient struct {
		// config stores application configuration
		config *config.Config
	}

	HourlyResponseData struct {
		Properties Properties `json:"properties"`
	}

	Properties struct {
		Periods []Period `json:"periods"`
	}

	Period struct {
		Number                       int    `json:"number"`
		Name                         string `json:"name"`
		StartTime                    string `json:"startTime"`
		EndTime                      string `json:"endTime"`
		IsDaytime                    bool   `json:"isDaytime"`
		Temperature                  int    `json:"temperature"`
		TemperatureUnit              string `json:"temperatureUnit"`
		ProbababilityOfPrecipitation struct {
			Unit  string `json:"unitCode"`
			Value int    `json:"value"`
		}
		Dewpoint struct {
			Unit  string  `json:"unitCode"`
			Value float64 `json:"value"`
		}
		RelativeHumidity struct {
			Unit  string `json:"unitCode"`
			Value int    `json:"value"`
		}

		WindSpeed        string `json:"windSpeed"`
		WindDirection    string `json:"windDirection"`
		Icon             string `json:"icon"`
		ShortForecast    string `json:"shortForecast"`
		DetailedForecast string `json:"detailedForecast"`
	}
)

// NewMailClient creates a new MailClient
func NewWeatherClient(cfg *config.Config) *WeatherClient {
	return &WeatherClient{
		config: cfg,
	}
}

func (w *WeatherClient) GetHourlyForecast(
	ctx echo.Context,
	x int,
	y int,
	office string,
) ([]Period, error) {
	ctx.Logger().Infof("Getting hourly forecast for x: %d, y: %d, office: %s", x, y, office)
	client := &http.Client{}
	url := fmt.Sprintf("%s/%s/%d,%d/forecast/hourly", w.config.Weather.Hostname, office, x, y)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		ctx.Logger().Errorf("unable to create hourly forecast request: %s", err)
		return nil, err
	}

	req.Header.Set("User-Agent", w.config.Weather.UserAgent)
	req.Header.Set("Accept", "application/json")
	ctx.Logger().Infof("Request: %v", req)
	resp, err := client.Do(req)
	ctx.Logger().Infof("Response: %v", resp)

	if err != nil {
		ctx.Logger().Errorf("unable to get hourly forecast: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	var data HourlyResponseData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	ctx.Logger().Infof("Data: %v", data)

	var periods []Period
	for _, period := range data.Properties.Periods {
		periods = append(periods, period)
		ctx.Logger().Infof("Period: %v", period)
	}

	return periods, nil
}
