package guilds

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/business/models/guildscache"
	"github.com/FlameInTheDark/rebot/business/models/guildsdb"
)

// Service is a guild service
type Service struct {
	guilds guildsdb.Querier
	cache  guildscache.GuildsCache

	logger *zap.Logger
}

// NewGuildsService creates a new guilds service
func NewGuildsService(db *sqlx.DB, rc *redis.Client, logger *zap.Logger) *Service {
	return &Service{
		guilds: guildsdb.New(db),
		cache:  guildscache.NewGuildsRedisCache(rc),
		logger: logger,
	}
}

// GetCommandPrefix get a guild command prefix from database or redis cache if available, else create one
func (g *Service) GetCommandPrefix(guildID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	cachedPrefix, err := g.cache.FindCommandPrefix(ctx, guildID)
	if err == nil {
		return cachedPrefix, nil
	}
	guild, err := g.guilds.Find(ctx, guildID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			guild, err = g.guilds.Create(ctx, guildID)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	err = g.cache.SetCommandPrefix(ctx, guildID, guild.CommandPrefix)
	if err != nil {
		g.logger.Error("unable cache guild prefix", zap.Error(err))
	}
	return guild.CommandPrefix, nil
}
