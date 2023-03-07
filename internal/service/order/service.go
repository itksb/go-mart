package order

import (
	"context"
	"errors"
	"github.com/itksb/go-mart/internal/domain"
	"github.com/itksb/go-mart/internal/service/order/luhn"
)

type Service struct {
	db domain.OrderRepositoryInterface
}

func NewOrderService(db domain.OrderRepositoryInterface) (*Service, error) {
	return &Service{db: db}, nil
}

func (o *Service) Create(ctx context.Context, orderID string, userID string) (*domain.Order, error) {
	if !luhn.Valid(orderID) {
		return nil, ErrOrderIncorrectOrderNumber
	}
	order, err := o.db.Create(ctx, orderID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrDomainDuplicateOrder) {
			orderFromDB, err2 := o.db.FindByID(ctx, orderID)
			if err2 != nil {
				return nil, err2
			}
			if orderFromDB.UserID == userID {
				return nil, ErrOrderAlreadyCreatedByCurUser
			}
			return nil, ErrOrderAlreadyCreatedByAnotherUser
		}
		return nil, err
	}

	return order, nil
}

func (o *Service) FindByID(ctx context.Context, ID string) (*domain.Order, error) {
	return o.db.FindByID(ctx, ID)
}

func (o *Service) FindAllByUserID(ctx context.Context, userID string) ([]*domain.Order, error) {
	return o.db.FindAllByUserID(ctx, userID)
}

func (o *Service) MakeAccrualForOrder(ctx context.Context, orderID string, status domain.OrderStatus, accrual float64) (*domain.Order, error) {
	return o.db.MakeAccrualForOrder(ctx, orderID, status, accrual)
}
