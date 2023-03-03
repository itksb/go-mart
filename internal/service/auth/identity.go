package auth

import "context"

type IdentityInterface interface {
	Create(ctx context.Context, login string, passHash string) (*User, error)
	FindOne(ctx context.Context, login string, userID string) (*User, error)
}
