package pgidentity

import (
	"context"
	"github.com/itksb/go-mart/internal/service/auth"
	"github.com/jmoiron/sqlx"
)

type IdentityPostgres struct {
	db *sqlx.DB
}

func NewIdentityPostgres(db *sqlx.DB) (*IdentityPostgres, error) {
	return &IdentityPostgres{db: db}, nil
}

func (i *IdentityPostgres) Create(ctx context.Context, params auth.IdentityParamsCreate) (*auth.User, error) {
	//TODO implement me
	panic("implement me")
}

func (i *IdentityPostgres) FindOne(ctx context.Context, params auth.IdentityParamsFindOne) (*auth.User, error) {
	//TODO implement me
	panic("implement me")
}
