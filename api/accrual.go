package api

// Взаимодействие с системой расчёта начислений баллов лояльности

type AccrualOrderRequest struct {
	Value string
}

type AccrualOrderResponse struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}
