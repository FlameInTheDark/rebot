package weather

import (
	"bytes"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/business/models/weatherdb"
	"github.com/FlameInTheDark/rebot/foundation/geonames"
	"github.com/FlameInTheDark/rebot/foundation/owm"
	"github.com/FlameInTheDark/rebot/foundation/wgen"
)

type Service struct {
	generator *wgen.Generator
	geonames  *geonames.Client
	owm       *owm.Client

	db weatherdb.Querier

	logger *zap.Logger
}

func NewService(generator *wgen.Generator, geo *geonames.Client, ow *owm.Client, db *sqlx.DB, logger *zap.Logger) *Service {
	return &Service{
		generator: generator,
		geonames:  geo,
		owm:       ow,
		db:        weatherdb.New(db),
		logger:    logger,
	}
}

func (s *Service) GetWeather(userId, location string) (*bytes.Buffer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var isOld bool
	if location == "" {
		var err error
		location, err = s.db.Find(ctx, userId)
		if err != nil {
			return nil, err
		}
		isOld = true
	}
	loc, err := s.geonames.FindOneLocation(ctx, location)
	if err != nil {
		s.logger.Debug("Find location error", zap.Error(err))
		return nil, err
	}
	lat, lng, err := loc.CoordinatesFloat64()
	if err != nil {
		s.logger.Debug("Getting coordinates error", zap.Error(err))
		return nil, err
	}
	forecast, err := s.owm.GetForecast(ctx, lat, lng, owm.ExcludeDaily+","+owm.ExcludeMinutely)
	if err != nil {
		s.logger.Debug("Getting forecast error", zap.Error(err))
		return nil, err
	}
	fd, err := convertHourly(forecast, loc.CountryName+", "+loc.Name)
	if err != nil {
		s.logger.Debug("Data convertation error", zap.Error(err))
		return nil, err
	}
	if !isOld {
		err = s.db.Insert(ctx, weatherdb.InsertParams{DiscordID: userId, Location: location})
		if err != nil {
			s.logger.Error("Cannot save weather location", zap.Error(err))
		}
	}
	return s.generator.Generate(fd)
}

func (s *Service) GetWeatherDaily(userId, location string) (*bytes.Buffer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var isOld bool
	if location == "" {
		var err error
		location, err = s.db.Find(ctx, userId)
		if err != nil {
			return nil, err
		}
		isOld = true
	}

	loc, err := s.geonames.FindOneLocation(ctx, location)
	if err != nil {
		s.logger.Debug("Find location error", zap.Error(err))
		return nil, err
	}
	lat, lng, err := loc.CoordinatesFloat64()
	if err != nil {
		s.logger.Debug("Getting coordinates error", zap.Error(err))
		return nil, err
	}
	forecast, err := s.owm.GetForecast(ctx, lat, lng, owm.ExcludeHourly+","+owm.ExcludeMinutely)
	if err != nil {
		s.logger.Debug("Getting forecast error", zap.Error(err))
		return nil, err
	}
	fd, err := convertDaily(forecast, loc.CountryName+", "+loc.Name)
	if err != nil {
		s.logger.Debug("Data convertation error", zap.Error(err))
		return nil, err
	}

	if !isOld {
		err = s.db.Insert(ctx, weatherdb.InsertParams{DiscordID: userId, Location: location})
		if err != nil {
			s.logger.Error("Cannot save weather location", zap.Error(err))
		}
	}
	return s.generator.GenerateDaily(fd)
}
