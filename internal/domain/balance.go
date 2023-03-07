package domain

import (
	"context"
)

type Balance struct {
	UserID  string
	Balance float64
}

type BalanceRepositoryInterface interface {
	FindByUserID(ctx context.Context, userID string) (*Balance, error)
	SumWithdrawnByUserID(ctx context.Context, userID string) (float64, error)
}
