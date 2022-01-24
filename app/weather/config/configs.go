package config

import (
	"sync"

	"github.com/google/uuid"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

var lock = &sync.Mutex{}

type Config struct {
	Http struct {
		Port int `env:"HTTP_PORT" env-default:"8080"`
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
	Influx struct {
		Host   string `env:"INFLUX_HOST" env-default:"influxdb:8086"`
		Token  string `env:"INFLUX_TOKEN" env-required:""`
		Org    string `env:"INFLUX_ORG" env-required:""`
		Bucket string `env:"INFLUX_BUCKET" env-required:""`
	}
	RabbitMQ struct {
		Host     string `env:"RABBIT_HOST" env-default:"mq"`
		Port     int    `env:"RABBIT_PORT" env-default:"5672"`
		User     string `env:"RABBIT_USER" env-default:"rabbitmq"`
		Password string `env:"RABBIT_PASS" env-default:"rabbitpass"`
	}
	Consul struct {
		Address     string `env:"CONSUL_ADDR" env-default:"consul:8500"`
		ServiceID   UUID   `env:"CONSUL_ID" env-default:""`
		ServiceName string `env:"CONSUL_NAME" env-default:"command"`
	}
	Location struct {
		GeonamesUsername string `env:"LOCATION_GEONAMES_USERNAME" env-default:""`
	}
	Forecast struct {
		OWMKey string `env:"OWM_API_KEY" env-default:""`
	}
	Weather struct {
		FontFile          string `env:"WEATHER_FONT_FILE" env-default:"lato.ttf"`
		IconsFile         string `env:"WEATHER_ICONS_FILE" env-default:"weathericons.ttf"`
		IconsBindingsFile string `env:"WEATHER_ICONS_BINDINGS_FILE" env-default:"icons_binds.json"`
	}
}

type UUID string

func (u *UUID) SetValue(s string) error {
	if s != "" {
		id, err := uuid.Parse(s)
		if err != nil {
			return err
		}
		*u = UUID(id.String())
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
