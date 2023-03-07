package domain

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

// Withdraw - Запрос на списание средств
type Withdraw struct {
	OrderID     string    `json:"order"`
	Amount      float64   `json:"sum"`
	UserID      string    `json:"-"`
	ProcessedAt time.Time `json:"processed_at"`
}

type WithdrawRepositoryInterface interface {
	Create(ctx context.Context, orderID string, amount float64, userID string) (*Withdraw, error)
	FindAllByUserID(ctx context.Context, userID string) ([]*Withdraw, error)
}

var ErrWithdrawNotEnoughBalance = errors.New("withdraw not enough balance")

func (w Withdraw) MarshalJSON() ([]byte, error) {
	type WithdrawAlias Withdraw
	return json.Marshal(&struct {
		WithdrawAlias
		ProcessedAt string `json:"processed_at"`
	}{
		WithdrawAlias: WithdrawAlias(w),
		ProcessedAt:   w.ProcessedAt.Format(time.RFC3339),
	})
}
