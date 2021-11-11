package guildscache

import "context"

type GuildsCache interface {
	FindCommandPrefix(ctx context.Context, guildId string) (string, error)
	SetCommandPrefix(ctx context.Context, guildId string, prefix string) error
}