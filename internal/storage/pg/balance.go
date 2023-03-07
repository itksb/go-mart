package pg

import (
	"context"
	"database/sql"
	"github.com/itksb/go-mart/internal/domain"
	"github.com/jmoiron/sqlx"
)

type BalanceRepositoryPg struct {
	db *sqlx.DB
}

func NewBalanceRepositoryPg(db *sqlx.DB) (*BalanceRepositoryPg, error) {
	return &BalanceRepositoryPg{db: db}, nil
}

func (r *BalanceRepositoryPg) FindByUserID(ctx context.Context, userID string) (*domain.Balance, error) {
	balance := domain.Balance{
		UserID:  userID,
		Balance: 0,
	}

	err := r.db.QueryRowxContext(
		ctx,
		"SELECT  balance FROM users WHERE id=$1",
		userID,
	).Scan(
		&balance.Balance,
	)

	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (r *BalanceRepositoryPg) SumWithdrawnByUserID(ctx context.Context, userID string) (float64, error) {
	var sum sql.NullFloat64
	err := r.db.QueryRowxContext(
		ctx,
		"select sum(amount) from withdraws where user_id = $1",
		userID,
	).Scan(&sum)

	if err != nil {
		return 0, err
	}

	if !sum.Valid {
		return 0, nil
	}

	return sum.Float64, nil

}
