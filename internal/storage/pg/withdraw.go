package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/itksb/go-mart/internal/domain"
	"github.com/jmoiron/sqlx"
	"strings"
)

type WithdrawRepositoryPg struct {
	db *sqlx.DB
}

func NewWithdrawRepositoryPg(db *sqlx.DB) (*WithdrawRepositoryPg, error) {
	return &WithdrawRepositoryPg{
		db: db,
	}, nil
}

func (w *WithdrawRepositoryPg) Create(ctx context.Context, orderID string, amount float64, userID string) (*domain.Withdraw, error) {
	withdraw := domain.Withdraw{OrderID: orderID, Amount: amount, UserID: userID}

	tx, err := w.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = tx.QueryRow(
		"UPDATE users SET balance = balance - $1 where id = $2",
		amount,
		userID,
	).Err()

	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(
		"INSERT INTO withdraws(order_id, amount, user_id) values($1, $2, $3) returning order_id, processed_at",
		orderID,
		amount,
		userID,
	).Scan(
		&withdraw.OrderID,
		&withdraw.ProcessedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates") {
			return nil, fmt.Errorf("order already exists. %w. %s", domain.ErrDomainDuplicateOrder, err.Error())
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &withdraw, nil
}

func (w *WithdrawRepositoryPg) FindAllByUserID(ctx context.Context, userID string) ([]*domain.Withdraw, error) {
	withdraws := make([]*domain.Withdraw, 0)

	rows, err := w.db.QueryxContext(
		ctx,
		"select  order_id, amount, user_id, processed_at from withdraws where user_id = $1 order by processed_at DESC",
		userID,
	)

	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		withdraw := &domain.Withdraw{}
		err = rows.Scan(
			&withdraw.OrderID,
			&withdraw.Amount,
			&withdraw.UserID,
			&withdraw.ProcessedAt,
		)
		if err == nil {
			withdraws = append(withdraws, withdraw)
		}
	}

	return withdraws, nil
}
