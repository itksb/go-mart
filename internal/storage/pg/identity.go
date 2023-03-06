package pg

import (
	"context"
	"github.com/itksb/go-mart/internal/service/auth"
	"github.com/jmoiron/sqlx"
)

type PgIdentityProvider struct {
}

func NewPgIdentityProvider(db sqlx.DB) (*PgIdentityProvider, error) {

}

func (PgIdentityProvider) Create(requestContext context.Context, params auth.IdentityParamsCreate) (*auth.User, error) {
	//TODO implement me
	panic("implement me")
}

func (PgIdentityProvider) FindOne(requestContext context.Context, params auth.IdentityParamsFindOne) (*auth.User, error) {
	//TODO implement me
	panic("implement me")
}
