package token

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("Token not found")

type Repository interface {
	GetToken(ctx context.Context, userId string) (TokenSchema, error)

	RemoveToken(ctx context.Context, schema TokenSchema) error

	SaveToken(ctx context.Context, schema TokenSchema) error
}
