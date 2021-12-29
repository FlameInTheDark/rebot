package service

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/business/service/weather"
	"github.com/FlameInTheDark/rebot/business/transport/commandst"
	"github.com/FlameInTheDark/rebot/foundation/geonames"
	"github.com/FlameInTheDark/rebot/foundation/owm"
	"github.com/FlameInTheDark/rebot/foundation/wgen"
)

type WeatherService struct {
	ws     *weather.Service
	cr     commandst.CommandsReceiver
	sess   *discordgo.Session
	logger *zap.Logger
}

func NewWeatherService(wg *wgen.Generator, fc *owm.Client, lc *geonames.Client, cr commandst.CommandsReceiver, sess *discordgo.Session, logger *zap.Logger) *WeatherService {
	return &WeatherService{ws: weather.NewService(wg, lc, fc, logger), cr: cr, sess: sess, logger: logger}
}

func (w *WeatherService) Run() error {
	w.cr.AddHandler(func(m commandst.CommandMessage) {
		data, err := w.ws.GetWeather(m.Message)
		if err != nil {
			w.logger.Debug("Cannot get weather data", zap.Error(err))
			return
		}
		_, err = w.sess.ChannelFileSend(m.ChannelID, fmt.Sprintf("weather_%d.png", time.Now().Unix()), data)
		if err != nil {
			w.logger.Debug("Cannot send weather picture", zap.String("channel-id", m.ChannelID), zap.String("user-id", m.UserID))
			return
		}
	})

	err := w.sess.Open()
	if err != nil {
		return err
	}

	return w.cr.Start("w")
}
