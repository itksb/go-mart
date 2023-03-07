package withdraw

import (
	"context"
	"github.com/itksb/go-mart/internal/domain"
	"github.com/itksb/go-mart/internal/service/order"
	"github.com/itksb/go-mart/internal/service/order/luhn"
)

type Service struct {
	db  domain.WithdrawRepositoryInterface
	bal domain.BalanceRepositoryInterface
}

func NewWithdrawService(db domain.WithdrawRepositoryInterface, bal domain.BalanceRepositoryInterface) (*Service, error) {
	return &Service{db: db, bal: bal}, nil
}

func (w *Service) Create(ctx context.Context, orderID string, amount float64, userID string) (*domain.Withdraw, error) {
	if !luhn.Valid(orderID) {
		return nil, order.ErrOrderIncorrectOrderNumber
	}
	userBalance, err := w.bal.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if userBalance.Balance < amount {
		return nil, domain.ErrWithdrawNotEnoughBalance
	}

	return w.db.Create(ctx, orderID, amount, userID)
}

func (w *Service) FindAllByUserID(ctx context.Context, userID string) ([]*domain.Withdraw, error) {
	return w.db.FindAllByUserID(ctx, userID)
}
