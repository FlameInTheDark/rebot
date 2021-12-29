package weather

import (
	"github.com/FlameInTheDark/rebot/foundation/owm"
	"github.com/FlameInTheDark/rebot/foundation/wgen"
	"github.com/pkg/errors"
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
