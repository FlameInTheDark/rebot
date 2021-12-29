package weather

import (
	"bytes"

	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/foundation/geonames"
	"github.com/FlameInTheDark/rebot/foundation/owm"
	"github.com/FlameInTheDark/rebot/foundation/wgen"
)

type Service struct {
	generator *wgen.Generator
	geonames  *geonames.Client
	owm       *owm.Client

	logger *zap.Logger
}

func NewService(generator *wgen.Generator, geo *geonames.Client, ow *owm.Client, logger *zap.Logger) *Service {
	return &Service{
		generator: generator,
		geonames:  geo,
		owm:       ow,
		logger: logger,
	}
}

func (s *Service) GetWeather(location string) (*bytes.Buffer, error) {
	loc, err := s.geonames.FindOneLocation(location)
	if err != nil {
		s.logger.Debug("Find location error", zap.Error(err))
		return nil, err
	}
	lat, lng, err := loc.CoordinatesFloat64()
	if err != nil {
		s.logger.Debug("Getting coordinates error", zap.Error(err))
		return nil, err
	}
	forecast, err := s.owm.GetForecast(lat, lng, owm.ExcludeDaily+","+owm.ExcludeMinutely)
	if err != nil {
		s.logger.Debug("Getting forecast error", zap.Error(err))
		return nil, err
	}
	fd, err := convertHourly(forecast, loc.CountryName+", "+loc.Name)
	if err != nil {
		s.logger.Debug("Data convertation error", zap.Error(err))
		return nil, err
	}
	return s.generator.Generate(fd)
}
