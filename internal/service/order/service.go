package order

import (
	"context"
	"errors"
	"github.com/itksb/go-mart/internal/domain"
)

type Order struct {
	db domain.OrderRepositoryInterface
}

func NewOrderService(db domain.OrderRepositoryInterface) *Order {
	return &Order{db: db}
}

func (o *Order) Create(ctx context.Context, orderID string, userID string) (*domain.Order, error) {
	order, err := o.db.Create(ctx, orderID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrDomainDuplicateOrder) {
			foundOrder, foundErr := o.db.FindByID(ctx, orderID)
			if foundErr != nil {
				return nil, foundErr
			}
			if foundOrder.UserID == userID {
				return nil, ErrAlreadyUploaded
			}
			return nil, ErrAlreadyUploadedByAnother
		}
		return nil, err
	}

	_ = o.rabbitmq.Publish(ctx, []byte(order.Number))
	return order, nil
}

func (o *Order) FindByNumber(ctx context.Context, number string) (*model.Order, error) {
	return o.orders.FindByNumber(ctx, number)
}

func (o *Order) FindAllByUserUUID(ctx context.Context, userUUID string) ([]*model.Order, error) {
	return o.orders.FindAllByUserUUID(ctx, userUUID)
}

func (o *Order) AccrueByNumber(ctx context.Context, number string, status model.OrderStatus, accrual float64) (*model.Order, error) {
	return o.orders.AccrueByNumber(ctx, number, status, accrual)
}
