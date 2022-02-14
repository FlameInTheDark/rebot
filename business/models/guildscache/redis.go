package guildscache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"time"
)

type GuildsRedisCache struct {
	client *redis.Client
}

func NewGuildsRedisCache(rc *redis.Client) *GuildsRedisCache {
	return &GuildsRedisCache{
		client: rc,
	}
}

func (c *GuildsRedisCache) FindCommandPrefix(ctx context.Context, guildID string) (string, error) {
	res := c.client.Get(ctx, fmt.Sprintf("guilds.%s.prefix", guildID))
	if err := res.Err(); errors.Is(err, redis.Nil) {
		return "", errors.Wrap(err, "prefix not found")
	}
	return res.Val(), nil
}

func (c *GuildsRedisCache) SetCommandPrefix(ctx context.Context, guildID string, prefix string) error {
	res := c.client.Set(ctx, fmt.Sprintf("guilds.%s.prefix", guildID), prefix, time.Hour)
	if err := res.Err(); err != nil {
		return err
	}
	return nil
}
