package config

import (
	"github.com/google/uuid"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"sync"
)

var lock = &sync.Mutex{}

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
	Consul struct {
		Address     string `env:"CONSUL_ADDR" env-default:"consul:8500"`
		ServiceID   UUID   `env:"CONSUL_ID" env-default:""`
		ServiceName string `env:"CONSUL_NAME" env-default:"api"`
	}
}

type UUID string

func (u *UUID) SetValue(s string) error {
	if s != "" {
		*u = UUID(s)
		return nil
	}
	*u = UUID(uuid.NewString())
	return nil
}

func (u UUID) String() string {
	return string(u)
}

var config *Config

//GetConfig load config from ENVs or .env file. If config is already loaded, return it.
func GetConfig() (*Config, error) {
	if config == nil {
		lock.Lock()
		defer lock.Unlock()

		if config == nil {
			var cfg Config
			if err := cleanenv.ReadEnv(&cfg); err != nil {
				return nil, errors.Wrap(err, "api: Read ENVs error")
			}
			config = &cfg
		}
	}
	return config, nil
}
