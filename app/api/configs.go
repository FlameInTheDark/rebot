package main

import (
	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"os"
)

type Config struct {
	Http struct {
		Port int `env:"HTTP_PORT,default=8080"`
	}
	Database struct {
		Host     string `env:"DATABASE_HOST,default=db"`
		Port     int    `env:"DATABASE_PORT,default=5432"`
		Username string `env:"DATABASE_USERNAME,default=postgres"`
		Password string `env:"DATABASE_PASS"`
	}
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.Wrap(err, "api: Read ENVs from .env file")
		}
	}

	if _, err := env.UnmarshalFromEnviron(&cfg); err != nil {
		return nil, errors.Wrap(err, "api: Unmarshal ENVs to environment structure")
	}
	return &cfg, nil
}
