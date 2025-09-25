package user

import (
	"context"
	"errors"
)

var ErrAlreadyExists = errors.New("This user already exists")
var ErrNotFound = errors.New("This user not found")

type Repository interface {
	Register(ctx context.Context, schema UserSchema, skipPassword bool) (UserSchema, error)

	GetUserByEmail(ctx context.Context, email string) (UserSchema, error)

	GetUserById(ctx context.Context, userId string) (UserSchema, error)

	AddImageId(ctx context.Context, schema UserSchema, imgId string) error
}
