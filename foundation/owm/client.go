package owm

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

const (
	UnitsStandard = "standard"
	UnitsMetric   = "metric"
	UnitsImperial = "imperial"

	ExcludeCurrent  = "current"
	ExcludeMinutely = "minutely"
	ExcludeHourly   = "hourly"
	ExcludeDaily    = "daily"

	LanguageEnglish = "en"

	EndpointOneCall = "https://api.openweathermap.org/data/2.5/onecall"
)

// Client open weather map client
type Client struct {
	apiKey   string
	units    string
	language string
}

// NewClient create new open weather map client
func NewClient(apiKey string, units string, language string) *Client {
	return &Client{apiKey: apiKey, units: units, language: language}
}

// GetForecast requests Forecast from OpenWeatherMap API
func (c *Client) GetForecast(ctx context.Context, lat, lng float64, exclude string) (*Forecast, error) {
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(
		"%s?lat=%f&lon=%f&exclude=%s&lang=%s&units=%s&appid=%s",
		EndpointOneCall,
		lat,
		lng,
		exclude,
		c.language,
		c.units,
		c.apiKey), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errBody ForecastError
		err = json.NewDecoder(resp.Body).Decode(&errBody)
		if err != nil {
			return nil, errors.New("forecast error")
		}
		return nil, errors.New(errBody.Message)
	}

	var data Forecast

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// GetHourlyForecast get hourly weather from the OpenWeatherMap service
func (c *Client) GetHourlyForecast(ctx context.Context, lat, lng float64) ([]HourlyForecast, error) {
	data, err := c.GetForecast(ctx, lat, lng, ExcludeMinutely+","+ExcludeDaily)
	if err != nil {
		return nil, err
	}
	return data.Hourly, nil
}

// GetDailyForecast get daily weather from the OpenWeatherMap service
func (c *Client) GetDailyForecast(ctx context.Context, lat, lng float64) ([]HourlyForecast, error) {
	data, err := c.GetForecast(ctx, lat, lng, ExcludeMinutely+","+ExcludeHourly)
	if err != nil {
		return nil, err
	}
	return data.Hourly, nil
}
