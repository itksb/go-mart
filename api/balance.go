package api

import "time"

// Получение текущего баланса пользователя

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn int     `json:"withdrawn"`
}

type BalanceWithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

// Получение информации о выводе средств

type BalanceWithdrawItem struct {
	Order       string    `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type BalanceWithdrawListResponse []BalanceWithdrawItem
