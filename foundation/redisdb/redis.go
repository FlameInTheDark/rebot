package redisdb

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

func NewConnection(host string, port int, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	status := client.Ping(ctx)
	if status.Err() != nil {
		return nil, errors.Wrap(status.Err(), "unable to connect to redis")
	}

	return client, nil
}
