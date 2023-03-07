package domain

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID         string      `json:"number"`
	Accrual    float64     `json:"accrual,omitempty"`
	Status     OrderStatus `json:"status"`
	UserID     string      `json:"-"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

type OrderRepositoryInterface interface {
	Create(ctx context.Context, ID string, userID string) (*Order, error)
	FindByID(ctx context.Context, ID string) (*Order, error)
	FindAllByUserID(ctx context.Context, ID string) ([]*Order, error)
	MakeAccrualForOrder(ctx context.Context, ID string, status OrderStatus, accrual float64) (*Order, error)
}

var ErrDomainDuplicateOrder = errors.New("duplicate key value")

func (o *Order) MarshalJSON() ([]byte, error) {
	type OrderAlias Order
	return json.Marshal(&struct {
		OrderAlias
		UploadedAt string `json:"uploaded_at"`
	}{
		OrderAlias: OrderAlias(*o),
		UploadedAt: o.UploadedAt.Format(time.RFC3339),
	})
}
