package main

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type Config struct {
	Http struct {
		Port int `env:"HTTP_PORT" env-default:"8080"`
	}
	Database struct {
		Host       string `env:"DATABASE_HOST" env-default:"db"`
		Port       int    `env:"DATABASE_PORT" env-default:"5432"`
		Database   string `env:"DATABASE_DBNAME" env-default:"postgres"`
		Username   string `env:"DATABASE_USERNAME" env-default:"postgres"`
		Password   string `env:"DATABASE_PASSWORD" env-default:"postgres"`
		DisableTLS bool   `env:"DATABASE_DISABLE_TLS" env-default:"true"`
		CetrPath   string `env:"DATABASE_CERT_PATH"`
	}
	Metrics struct {
		Host   string `env:"METRICS_HOST" env-default:"influxdb"`
		Port   int    `env:"METRICS_PORT" env-default:"8086"`
		Bucket string `env:"METRICS_BUCKET" env-default:"rebot"`
		Token  string `env:"METRICS_TOKEN"`
	}
}

func GetConfig() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, errors.Wrap(err, "api: Read ENVs error")
	}
	return &cfg, nil
}
