package pg

import (
	"context"
	"errors"
	"fmt"
	"github.com/itksb/go-mart/internal/service/auth"
	"github.com/jmoiron/sqlx"
	"strings"
)

type IdentityProviderPg struct {
	db *sqlx.DB
}

func NewIdentityProviderPg(db *sqlx.DB) (*IdentityProviderPg, error) {
	return &IdentityProviderPg{
		db: db,
	}, nil
}

func (i *IdentityProviderPg) Create(requestContext context.Context, params auth.IdentityParamsCreate) (*auth.User, error) {
	authUser := auth.User{
		Login:        params.Login,
		PasswordHash: params.PassHash,
	}

	q := i.db.QueryRowxContext(
		requestContext,
		"INSERT INTO public.users(login, password_hash) VALUES ($1, $2) RETURNING id, created_at",
		params.Login,
		params.PassHash,
	)

	err := q.Scan(&authUser.ID, &authUser.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, fmt.Errorf("login already exists. %w. %s", auth.ErrDuplicateKeyValue, err.Error())
		}
		return nil, err
	}

	return &authUser, nil
}

func (i *IdentityProviderPg) FindOne(requestContext context.Context, params auth.IdentityParamsFindOne) (*auth.User, error) {
	if params.UserID == "" && params.Login == "" {
		return nil, errors.New("FindOne error. Params all empty")
	}

	authUser := auth.User{
		ID:    params.UserID,
		Login: params.Login,
	}

	var row *sqlx.Row

	if params.UserID != "" {
		row = i.db.QueryRowxContext(
			requestContext,
			"SELECT id, login, password_hash, created_at FROM  public.users WHERE id=$1",
			params.UserID,
		)
	} else if params.UserID != "" && params.Login != "" {
		row = i.db.QueryRowxContext(
			requestContext,
			"SELECT id, login, password_hash, created_at FROM  public.users WHERE id=$1 AND login=$2",
			params.UserID,
			params.Login,
		)
	} else {
		row = i.db.QueryRowxContext(
			requestContext,
			"SELECT id, login, password_hash, created_at FROM  public.users WHERE login=$1",
			params.Login,
		)
	}

	err := row.Scan(&authUser.ID, &authUser.Login, &authUser.PasswordHash, &authUser.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &authUser, nil

}
