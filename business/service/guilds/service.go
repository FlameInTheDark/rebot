package guilds

import (
	"context"
	"github.com/FlameInTheDark/rebot/business/models/guildscache"
	"github.com/FlameInTheDark/rebot/business/models/guildsdb"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type GuildsService struct {
	guilds guildsdb.Querier
	cache  guildscache.GuildsCache

	logger *zap.Logger
}

func NewGuildsService(db *sqlx.DB, rc *redis.Client, logger *zap.Logger) *GuildsService {
	return &GuildsService{
		guilds: guildsdb.New(db),
		cache:  guildscache.NewGuildsRedisCache(rc),
		logger: logger,
	}
}

func (g *GuildsService) GetCommandPrefix(guildId string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	cachedPrefix, err := g.cache.FindCommandPrefix(ctx, guildId)
	if err == nil {
		return cachedPrefix, nil
	}
	guild, err := g.guilds.Find(ctx, guildId)
	if err != nil {
		return "", err
	}
	err = g.cache.SetCommandPrefix(ctx, guildId, guild.CommandPrefix)
	if err != nil {
		g.logger.Error("unable cache guild prefix", zap.Error(err))
	}
	return guild.CommandPrefix, nil
}