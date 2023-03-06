package auth

import (
	"context"
	"errors"
)

type IdentityProviderInterface interface {
	Create(requestContext context.Context, params IdentityParamsCreate) (*User, error)
	FindOne(requestContext context.Context, params IdentityParamsFindOne) (*User, error)
}

type IdentityParamsCreate struct{ Login, PassHash string }
type IdentityParamsFindOne struct{ Login, UserID string }

var ErrDuplicateKeyValue = errors.New("duplicate key value")
