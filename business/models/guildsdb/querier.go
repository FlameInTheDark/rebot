// Code generated by sqlc. DO NOT EDIT.

package guildsdb

import (
	"context"
)

type Querier interface {
	Create(ctx context.Context, discordID string) (Guild, error)
	Find(ctx context.Context, discordID string) (Guild, error)
}

var _ Querier = (*Queries)(nil)
