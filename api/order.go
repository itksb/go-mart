package api

import "time"

// OrderLoadRequest - Загрузка номера заказа
type OrderLoadRequest struct {
	Value string
}

// Получение списка загруженных номеров заказов

type OrderItem struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type OrderListResponse []OrderItem
