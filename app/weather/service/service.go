package service

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/app/weather/metrics"
	"github.com/FlameInTheDark/rebot/business/service/weather"
	"github.com/FlameInTheDark/rebot/business/transport/commandst"
	"github.com/FlameInTheDark/rebot/foundation/geonames"
	"github.com/FlameInTheDark/rebot/foundation/metricsclients"
	"github.com/FlameInTheDark/rebot/foundation/owm"
	"github.com/FlameInTheDark/rebot/foundation/wgen"
)

type WeatherService struct {
	ws      *weather.Service
	cr      commandst.CommandsReceiver
	sess    *discordgo.Session
	metrics *metrics.Metrics
	logger  *zap.Logger
}

func NewWeatherService(
	wg *wgen.Generator,
	fc *owm.Client,
	lc *geonames.Client,
	cr commandst.CommandsReceiver,
	sess *discordgo.Session,
	db *sqlx.DB,
	mc *metricsclients.InfluxMetrics,
	logger *zap.Logger,
) *WeatherService {
	return &WeatherService{
		ws:      weather.NewService(wg, lc, fc, db, logger),
		cr:      cr,
		sess:    sess,
		metrics: metrics.NewMetrics(mc),
		logger:  logger,
	}
}

func (w *WeatherService) Run() error {
	w.cr.SetErrorMetrics(w.metrics)

	w.cr.AddHandler("weather", func(m commandst.CommandMessage) error {
		data, err := w.ws.GetWeather(m.UserID, m.Message)
		if err != nil {
			return errors.Wrap(err, "Cannot get weather data")
		}

		defer func() {
			data.Reset()
		}()

		_, err = w.sess.ChannelFileSend(m.ChannelID, fmt.Sprintf("weather_%d.png", time.Now().Unix()), data)
		if err != nil {
			return errors.Wrap(err, "Cannot send weather picture")
		}
		return nil
	})

	w.cr.AddHandler("wweather", func(m commandst.CommandMessage) error {
		data, err := w.ws.GetWeatherDaily(m.UserID, m.Message)
		if err != nil {
			return errors.Wrap(err, "Cannot get weather data")
		}

		defer func() {
			data.Reset()
		}()

		_, err = w.sess.ChannelFileSend(m.ChannelID, fmt.Sprintf("weather_%d.png", time.Now().Unix()), data)
		if err != nil {
			return errors.Wrap(err, "Cannot send weather picture")
		}
		return nil
	})

	err := w.sess.Open()
	if err != nil {
		return err
	}

	err = w.cr.Start("weather")
	if err != nil {
		return err
	}

	err = w.cr.Start("wweather")
	if err != nil {
		return err
	}

	return nil
}
