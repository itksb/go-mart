package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/itksb/go-mart/internal/domain"
	"github.com/jmoiron/sqlx"
	"strings"
)

type OrderRepositoryPg struct {
	db *sqlx.DB
}

func NewOrderRepositoryPg(db *sqlx.DB) (*OrderRepositoryPg, error) {
	return &OrderRepositoryPg{db: db}, nil
}

func (r *OrderRepositoryPg) Create(ctx context.Context, ID string, userID string) (*domain.Order, error) {
	order := domain.Order{
		ID:      ID,
		Accrual: 0,
		Status:  domain.OrderStatusNew,
		UserID:  userID,
	}

	err := r.db.QueryRowxContext(
		ctx,
		"INSERT INTO orders(id, status, user_id)  VALUES ($1, $2, $3) RETURNING id, accrual, status, user_id, uploaded_at",
		ID,
		domain.OrderStatusNew,
		userID,
	).Scan(
		&order.ID,
		&order.Accrual,
		&order.Status,
		&order.UserID,
		&order.UploadedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, fmt.Errorf("order already exists. %w. %s", domain.ErrDomainDuplicateOrder, err.Error())
		}
		return nil, err
	}

	return &order, nil

}

func (r *OrderRepositoryPg) FindByID(ctx context.Context, ID string) (*domain.Order, error) {
	order := domain.Order{ID: ID}
	err := r.db.QueryRowxContext(
		ctx,
		"SELECT id, accrual, status, user_id, uploaded_at FROM orders WHERE id=$1",
		ID,
	).Scan(
		&order.ID,
		&order.Accrual,
		&order.Status,
		&order.UserID,
		&order.UploadedAt,
	)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepositoryPg) FindAllByUserID(ctx context.Context, ID string) ([]*domain.Order, error) {
	orders := make([]*domain.Order, 0)
	rows, err := r.db.QueryxContext(
		ctx,
		"SELECT id, accrual, status, user_id, uploaded_at FROM orders WHERE user_id=$1 ORDER BY uploaded_at DESC",
		ID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		order := &domain.Order{}
		err = rows.Scan(
			&order.ID,
			&order.Accrual,
			&order.Status,
			&order.UserID,
			&order.UploadedAt,
		)
		if err == nil {
			orders = append(orders, order)
		}
	}

	return orders, nil

}

func (r *OrderRepositoryPg) MakeAccrualForOrder(ctx context.Context, ID string, status domain.OrderStatus, accrual float64) (*domain.Order, error) {
	order := domain.Order{
		ID:      ID,
		Status:  status,
		Accrual: accrual,
	}

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(
		"UPDATE orders SET status = $1, accrual = $2 WHERE id = $3 RETURNING id, accrual, status, user_id, uploaded_at",
		status,
		accrual,
		ID,
	).Scan(
		&order.ID,
		&order.Accrual,
		&order.Status,
		&order.UserID,
		&order.UploadedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("%s. %w", err.Error(), tx.Rollback())
	}

	_, err = tx.Exec(
		"update users set balance = balance + $1 where id = $2",
		accrual,
		order.UserID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s. %w", err.Error(), tx.Rollback())
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s. %w", err.Error(), tx.Rollback())
	}

	return &order, nil
}
