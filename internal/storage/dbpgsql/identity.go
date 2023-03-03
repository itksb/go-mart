package dbpgsql

import (
	"context"
	"github.com/itksb/go-mart/internal/service/auth"
)

type IdentityPostgres struct {
	db *PgStorage
}

func NewIdentityPostgres(db *PgStorage) (*IdentityPostgres, error) {
	return &IdentityPostgres{db: db}, nil
}

func (IdentityPostgres) Create(ctx context.Context, login string, passHash string) (*auth.User, error) {
	//TODO implement me
	panic("implement me")
}

func (IdentityPostgres) FindOne(ctx context.Context, login string, userID string) (*auth.User, error) {
	//TODO implement me
	panic("implement me")
}
