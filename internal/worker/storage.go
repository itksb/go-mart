package worker

import (
	"context"
	"github.com/itksb/go-mart/internal/domain"
	"github.com/jmoiron/sqlx"
)

type WorkerStorage struct {
	db *sqlx.DB
}

func NewWorkerStorage(db *sqlx.DB) (*WorkerStorage, error) {
	return &WorkerStorage{
		db: db,
	}, nil
}

func (w *WorkerStorage) LoadOrdersToWork(ctx context.Context) ([]*domain.Order, error) {
	orders := make([]*domain.Order, 0)
	rows, err := w.db.QueryxContext(
		ctx,
		"SELECT id,  status, user_id FROM orders WHERE status IN( 'NEW','PROCESSING')",
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
		order := &domain.Order{}
		err = rows.Scan(
			&order.ID,
			&order.Status,
			&order.UserID,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)

	}

	return orders, nil
}

func (w *WorkerStorage) UpdateOrderStatus(accrualResponse AccrualResponse) error {
	if accrualResponse.Status == domain.OrderStatusProcessed {
		var balance float32
		err := w.db.QueryRow(
			"SELECT balance FROM users WHERE id = $1",
			accrualResponse.UserID,
		).Scan(&balance)
		if err != nil {
			return err
		}

		tx, err := w.db.Begin()
		if err != nil {
			return err
		}
		stmtUser, err := tx.Prepare("UPDATE users SET balance=$1 WHERE id = $2")
		if err != nil {
			return err
		}
		defer stmtUser.Close()

		stmtOrder, err := tx.Prepare("UPDATE orders SET status='PROCESSED', accrual=$1 WHERE id=$2")
		if err != nil {
			return err
		}
		defer stmtOrder.Close()

		if _, err = stmtUser.Exec(balance+accrualResponse.Accrual, accrualResponse.UserID); err != nil {
			if err = tx.Rollback(); err != nil {
				return err
			}

		}
		if _, err = stmtOrder.Exec(accrualResponse.Accrual, accrualResponse.Order); err != nil {
			if err = tx.Rollback(); err != nil {
				return err
			}
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		return nil
	}

	_, err := w.db.Exec("UPDATE orders SET status=$1 WHERE id=$2", accrualResponse.Status, accrualResponse.Order)
	if err != nil {
		return err
	}
	return nil
}
