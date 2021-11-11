package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"sync"
)

var lock = &sync.Mutex{}

type Config struct {
	Http struct {
		Port string `env:"HTTP_PORT" env-default:"8080"`
	}
	Discord struct {
		Token string `env:"DISCORD_TOKEN"`
	}
	Database struct {
		Host       string `env:"DATABASE_HOST" env-default:"db"`
		Port       int    `env:"DATABASE_PORT" env-default:"5432"`
		Database   string `env:"DATABASE_DBNAME" env-default:"postgres"`
		Username   string `env:"DATABASE_USERNAME" env-default:"postgres"`
		Password   string `env:"DATABASE_PASSWORD" env-default:"postgres"`
		DisableTLS bool   `env:"DATABASE_DISABLE_TLS" env-default:"true"`
		CertPath   string `env:"DATABASE_CERT_PATH"`
	}
	Redis struct {
		Host     string `env:"REDIS_HOST" env-default:"redis"`
		Port     int    `env:"REDIS_PORT" env-default:"6379"`
		Password string `env:"REDIS_PASSWORD" env-default:"redispassword"`
		Database int    `env:"REDIS_DATABASE" env-default:"0"`
	}
	Metrics struct {
		Host   string `env:"METRICS_HOST" env-default:"influxdb"`
		Port   int    `env:"METRICS_PORT" env-default:"8086"`
		Bucket string `env:"METRICS_BUCKET" env-default:"rebot"`
		Token  string `env:"METRICS_TOKEN"`
	}
	RabbitMQ struct {
		Host     string `env:"RABBIT_HOST" env-default:"mq"`
		Port     int    `env:"RABBIT_PORT" env-default:"5672"`
		User     string `env:"RABBIT_USER" env-default:"rabbitmq"`
		Password string `env:"RABBIT_PASS" env-default:"rabbitpass"`
	}
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
