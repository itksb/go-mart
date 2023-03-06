package auth

import "context"

type IdentityProviderInterface interface {
	Create(requestContext context.Context, params IdentityParamsCreate) (*User, error)
	FindOne(requestContext context.Context, params IdentityParamsFindOne) (*User, error)
}

type IdentityParamsCreate struct{ login, passHash string }
type IdentityParamsFindOne struct{ login, userID string }
