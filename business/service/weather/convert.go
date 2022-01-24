package weather

import (
	"github.com/pkg/errors"

	"github.com/FlameInTheDark/rebot/foundation/owm"
	"github.com/FlameInTheDark/rebot/foundation/wgen"
)

func convertHourly(f *owm.Forecast, location string) (*wgen.ForecastData, error) {
	if len(f.Hourly) < 9 {
		return nil, errors.New("not enough data")
	}

	var fd wgen.ForecastData

	fd.Location = location
	fd.Forecast = append(fd.Forecast, wgen.ForecastRow{
		IconCode:    f.Current.Weather[0].Icon,
		Temperature: f.Current.Temperature,
		Humidity:    f.Current.Humidity,
		Clouds:      f.Current.Clouds,
		Time:        f.Current.GetCurrentTime(),
	})
	fd.Forecast = append(fd.Forecast, wgen.ForecastRow{
		IconCode:    f.Hourly[2].Weather[0].Icon,
		Temperature: f.Hourly[2].Temperature,
		Humidity:    f.Hourly[2].Humidity,
		Clouds:      f.Hourly[2].Clouds,
		Time:        f.Hourly[2].GetForecastTime(),
	})
	fd.Forecast = append(fd.Forecast, wgen.ForecastRow{
		IconCode:    f.Hourly[4].Weather[0].Icon,
		Temperature: f.Hourly[4].Temperature,
		Humidity:    f.Hourly[4].Humidity,
		Clouds:      f.Hourly[4].Clouds,
		Time:        f.Hourly[4].GetForecastTime(),
	})
	fd.Forecast = append(fd.Forecast, wgen.ForecastRow{
		IconCode:    f.Hourly[6].Weather[0].Icon,
		Temperature: f.Hourly[6].Temperature,
		Humidity:    f.Hourly[6].Humidity,
		Clouds:      f.Hourly[6].Clouds,
		Time:        f.Hourly[6].GetForecastTime(),
	})
	fd.Forecast = append(fd.Forecast, wgen.ForecastRow{
		IconCode:    f.Hourly[8].Weather[0].Icon,
		Temperature: f.Hourly[8].Temperature,
		Humidity:    f.Hourly[8].Humidity,
		Clouds:      f.Hourly[8].Clouds,
		Time:        f.Hourly[8].GetForecastTime(),
	})

	return &fd, nil
}

func convertDaily(f *owm.Forecast, location string) (*wgen.ForecastData, error) {
	if len(f.Daily) < 5 {
		return nil, errors.New("not enough data")
	}

	var fd wgen.ForecastData

	fd.Location = location
	fd.Forecast = append(fd.Forecast, wgen.ForecastRow{
		IconCode:    f.Current.Weather[0].Icon,
		Temperature: f.Current.Temperature,
		Humidity:    f.Current.Humidity,
		Clouds:      f.Current.Clouds,
		Time:        f.Current.GetCurrentTime(),
	})
	fd.Forecast = append(fd.Forecast, wgen.ForecastRow{
		IconCode:    f.Daily[1].Weather[0].Icon,
		Temperature: f.Daily[1].Temperature.Eve,
		Max:         f.Daily[1].Temperature.Max,
		Min:         f.Daily[1].Temperature.Min,
		Humidity:    f.Daily[1].Humidity,
		Clouds:      f.Daily[1].Clouds,
		Time:        f.Daily[1].GetForecastTime(),
	})
	fd.Forecast = append(fd.Forecast, wgen.ForecastRow{
		IconCode:    f.Daily[2].Weather[0].Icon,
		Temperature: f.Daily[2].Temperature.Eve,
		Max:         f.Daily[2].Temperature.Max,
		Min:         f.Daily[2].Temperature.Min,
		Humidity:    f.Daily[2].Humidity,
		Clouds:      f.Daily[2].Clouds,
		Time:        f.Daily[2].GetForecastTime(),
	})
	fd.Forecast = append(fd.Forecast, wgen.ForecastRow{
		IconCode:    f.Daily[3].Weather[0].Icon,
		Temperature: f.Daily[3].Temperature.Eve,
		Max:         f.Daily[3].Temperature.Max,
		Min:         f.Daily[3].Temperature.Min,
		Humidity:    f.Daily[3].Humidity,
		Clouds:      f.Daily[3].Clouds,
		Time:        f.Daily[3].GetForecastTime(),
	})
	fd.Forecast = append(fd.Forecast, wgen.ForecastRow{
		IconCode:    f.Daily[4].Weather[0].Icon,
		Temperature: f.Daily[4].Temperature.Eve,
		Max:         f.Daily[4].Temperature.Max,
		Min:         f.Daily[4].Temperature.Min,
		Humidity:    f.Daily[4].Humidity,
		Clouds:      f.Daily[4].Clouds,
		Time:        f.Daily[4].GetForecastTime(),
	})

	return &fd, nil
}
