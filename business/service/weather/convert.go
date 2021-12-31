package weather

import (
	"github.com/pkg/errors"

	"github.com/FlameInTheDark/rebot/foundation/owm"
	"github.com/FlameInTheDark/rebot/foundation/wgen"
)

func convertHourly(f *owm.Forecast, location string) (*wgen.ForecastData, error) {
	if len(f.Hourly) < 8 {
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
		IconCode:    f.Hourly[1].Weather[0].Icon,
		Temperature: f.Hourly[1].Temperature,
		Humidity:    f.Hourly[1].Humidity,
		Clouds:      f.Hourly[1].Clouds,
		Time:        f.Hourly[1].GetForecastTime(),
	})
	fd.Forecast = append(fd.Forecast, wgen.ForecastRow{
		IconCode:    f.Hourly[3].Weather[0].Icon,
		Temperature: f.Hourly[3].Temperature,
		Humidity:    f.Hourly[3].Humidity,
		Clouds:      f.Hourly[3].Clouds,
		Time:        f.Hourly[3].GetForecastTime(),
	})
	fd.Forecast = append(fd.Forecast, wgen.ForecastRow{
		IconCode:    f.Hourly[5].Weather[0].Icon,
		Temperature: f.Hourly[5].Temperature,
		Humidity:    f.Hourly[5].Humidity,
		Clouds:      f.Hourly[5].Clouds,
		Time:        f.Hourly[5].GetForecastTime(),
	})
	fd.Forecast = append(fd.Forecast, wgen.ForecastRow{
		IconCode:    f.Hourly[7].Weather[0].Icon,
		Temperature: f.Hourly[7].Temperature,
		Humidity:    f.Hourly[7].Humidity,
		Clouds:      f.Hourly[7].Clouds,
		Time:        f.Hourly[7].GetForecastTime(),
	})

	return &fd, nil
}

func convertDaily(f *owm.Forecast, location string) (*wgen.ForecastData, error) {
	if len(f.Daily) < 4 {
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
		IconCode:    f.Daily[0].Weather[0].Icon,
		Temperature: f.Daily[0].Temperature.Max,
		Humidity:    f.Daily[0].Humidity,
		Clouds:      f.Daily[0].Clouds,
		Time:        f.Daily[0].GetForecastTime(),
	})
	fd.Forecast = append(fd.Forecast, wgen.ForecastRow{
		IconCode:    f.Daily[1].Weather[0].Icon,
		Temperature: f.Daily[1].Temperature.Max,
		Humidity:    f.Daily[1].Humidity,
		Clouds:      f.Daily[1].Clouds,
		Time:        f.Daily[1].GetForecastTime(),
	})
	fd.Forecast = append(fd.Forecast, wgen.ForecastRow{
		IconCode:    f.Daily[2].Weather[0].Icon,
		Temperature: f.Daily[2].Temperature.Max,
		Humidity:    f.Daily[2].Humidity,
		Clouds:      f.Daily[2].Clouds,
		Time:        f.Daily[2].GetForecastTime(),
	})
	fd.Forecast = append(fd.Forecast, wgen.ForecastRow{
		IconCode:    f.Daily[3].Weather[0].Icon,
		Temperature: f.Daily[3].Temperature.Max,
		Humidity:    f.Daily[3].Humidity,
		Clouds:      f.Daily[3].Clouds,
		Time:        f.Daily[3].GetForecastTime(),
	})

	return &fd, nil
}
