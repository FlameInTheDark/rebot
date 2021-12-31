// Code generated by sqlc. DO NOT EDIT.

package weatherdb

import (
	"context"
)

type Querier interface {
	Create(ctx context.Context, arg CreateParams) (Weather, error)
	Find(ctx context.Context, discordID string) (string, error)
	Insert(ctx context.Context, arg InsertParams) error
	Update(ctx context.Context, location string) error
}

var _ Querier = (*Queries)(nil)
