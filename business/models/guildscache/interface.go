package guildscache

import "context"

// GuildsCache is a caching interface
type GuildsCache interface {
	FindCommandPrefix(ctx context.Context, guildID string) (string, error)
	SetCommandPrefix(ctx context.Context, guildID string, prefix string) error
}
