package owm

import "time"

type ForecastError struct {
	Message string `json:"message"`
	Cod     string `json:"cod"`
}

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

func (f *CurrentForecast) GetCurrentTime() time.Time {
	return time.Unix(f.CurrentTime, 0)
}

func (f *CurrentForecast) GetSunriseTime() time.Time {
	return time.Unix(f.SunriseTime, 0)
}

func (f *CurrentForecast) GetSunsetTime() time.Time {
	return time.Unix(f.SunsetTime, 0)
}

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

func (f *HourlyForecast) GetForecastTime() time.Time {
	return time.Unix(f.ForecastTime, 0)
}

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

func (d *DailyForecast) GetForecastTime() time.Time {
	return time.Unix(d.ForecastTime, 0)
}

func (d *DailyForecast) GetSunriseTime() time.Time {
	return time.Unix(d.SunriseTime, 0)
}

func (d *DailyForecast) GetSunsetTime() time.Time {
	return time.Unix(d.SunsetTime, 0)
}

func (d *DailyForecast) GetMoonriseTime() time.Time {
	return time.Unix(d.MoonriseTime, 0)
}

func (d *DailyForecast) GetMoonsetTime() time.Time {
	return time.Unix(d.MoonsetTime, 0)
}

type DailyTemperature struct {
	Day   float64 `json:"day"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"`
	Morn  float64 `json:"morn"`
}

type DailyFeelsLike struct {
	Day   float64 `json:"day"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"`
	Morn  float64 `json:"morn"`
}

type Weather struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Precipitation struct {
	H float64 `json:"1h"`
}

type Minutely struct {
	Dt            int `json:"dt"`
	Precipitation int `json:"precipitation"`
}
