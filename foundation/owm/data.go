package owm

import "time"

// ForecastError contains API error message
type ForecastError struct {
	Message string `json:"message"`
	Cod     string `json:"cod"`
}

// Forecast is a requested weather data
type Forecast struct {
	Lat            float64          `json:"lat"`
	Lon            float64          `json:"lon"`
	TimeZone       string           `json:"timezone"`
	TimeZoneOffset int              `json:"timezone_offset"`
	Current        CurrentForecast  `json:"current"`
	Minutely       []Minutely       `json:"minutely"`
	Hourly         []HourlyForecast `json:"hourly"`
	Daily          []DailyForecast  `json:"daily"`
}

// CurrentForecast contains the forecast for the current day
type CurrentForecast struct {
	CurrentTime int64         `json:"dt"`
	SunriseTime int64         `json:"sunrise"`
	SunsetTime  int64         `json:"sunset"`
	Temperature float64       `json:"temp"`
	FeelsLike   float64       `json:"feels_like"`
	Pressure    int           `json:"pressure"`
	Humidity    int           `json:"humidity"`
	DewPoint    float64       `json:"dew_point"`
	UVIndex     float64       `json:"uvi"`
	Clouds      int           `json:"clouds"`
	Visibility  int           `json:"visibility"`
	WindSpeed   float64       `json:"wind_speed"`
	WindGust    float64       `json:"wind_gust"`
	WindDeg     int           `json:"wind_deg"`
	Rain        Precipitation `json:"rain"`
	Snow        Precipitation `json:"snow"`
	Weather     []Weather     `json:"weather"`
}

// GetCurrentTime returns current time
func (f *CurrentForecast) GetCurrentTime() time.Time {
	return time.Unix(f.CurrentTime, 0)
}

// GetSunriseTime returns sunrise time
func (f *CurrentForecast) GetSunriseTime() time.Time {
	return time.Unix(f.SunriseTime, 0)
}

// GetSunsetTime returns sunset time
func (f *CurrentForecast) GetSunsetTime() time.Time {
	return time.Unix(f.SunsetTime, 0)
}

// HourlyForecast contains hourly weather forecast
type HourlyForecast struct {
	ForecastTime int64         `json:"dt"`
	Temperature  float64       `json:"temp"`
	FeelsLike    float64       `json:"feels_like"`
	Pressure     int           `json:"pressure"`
	Humidity     int           `json:"humidity"`
	DewPoint     float64       `json:"dew_point"`
	UVIndex      float64       `json:"uvi"`
	Clouds       int           `json:"clouds"`
	Visibility   int           `json:"visibility"`
	WindSpeed    float64       `json:"wind_speed"`
	WindGust     float64       `json:"wind_gust"`
	WindDeg      int           `json:"wind_deg"`
	Rain         Precipitation `json:"rain"`
	Snow         Precipitation `json:"snow"`
	Weather      []Weather     `json:"weather"`
}

// GetForecastTime returns forecast time
func (f *HourlyForecast) GetForecastTime() time.Time {
	return time.Unix(f.ForecastTime, 0)
}

// DailyForecast contains daily weather Forecast
type DailyForecast struct {
	ForecastTime int64            `json:"dt"`
	SunriseTime  int64            `json:"sunrise"`
	SunsetTime   int64            `json:"sunset"`
	MoonriseTime int64            `json:"moonrise"`
	MoonsetTime  int64            `json:"moonset"`
	MoonPhase    float64          `json:"moon_phase"`
	Temperature  DailyTemperature `json:"temp"`
	FeelsLike    DailyFeelsLike   `json:"feels_like"`
	Pressure     int              `json:"pressure"`
	Humidity     int              `json:"humidity"`
	DewPoint     float64          `json:"dew_point"`
	UVIndex      float64          `json:"uvi"`
	Clouds       int              `json:"clouds"`
	Visibility   int              `json:"visibility"`
	WindSpeed    float64          `json:"wind_speed"`
	WindGust     float64          `json:"wind_gust"`
	WindDeg      int              `json:"wind_deg"`
	Rain         float64          `json:"rain"`
	Weather      []Weather        `json:"weather"`
	POP          float64          `json:"pop"`
}

// GetForecastTime returns forecast time for the forecast day
func (d *DailyForecast) GetForecastTime() time.Time {
	return time.Unix(d.ForecastTime, 0)
}

// GetSunriseTime return sunrise time for the forecast day
func (d *DailyForecast) GetSunriseTime() time.Time {
	return time.Unix(d.SunriseTime, 0)
}

// GetSunsetTime returns sunset time for the forecast day
func (d *DailyForecast) GetSunsetTime() time.Time {
	return time.Unix(d.SunsetTime, 0)
}

// GetMoonriseTime returns moonrise time for the forecast day
func (d *DailyForecast) GetMoonriseTime() time.Time {
	return time.Unix(d.MoonriseTime, 0)
}

// GetMoonsetTime returns moonset time for the forecast day
func (d *DailyForecast) GetMoonsetTime() time.Time {
	return time.Unix(d.MoonsetTime, 0)
}

// DailyTemperature contains temperature data for the daily forecast
type DailyTemperature struct {
	Day   float64 `json:"day"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"`
	Morn  float64 `json:"morn"`
}

// DailyFeelsLike returns temperature that feels like for the forecast day
type DailyFeelsLike struct {
	Day   float64 `json:"day"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"`
	Morn  float64 `json:"morn"`
}

// Weather contains weather data
type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// Precipitation ...
type Precipitation struct {
	H float64 `json:"1h"`
}

// Minutely ...
type Minutely struct {
	Dt            int `json:"dt"`
	Precipitation int `json:"precipitation"`
}
