package handler

import (
	"github.com/itksb/go-mart/internal/config"
	"github.com/itksb/go-mart/internal/service/auth"
	"github.com/itksb/go-mart/internal/service/balance"
	"github.com/itksb/go-mart/internal/service/order"
	"github.com/itksb/go-mart/internal/service/withdraw"
	"github.com/itksb/go-mart/pkg/logger"
)

type Handler struct {
	logger          logger.Interface
	config          config.Config
	auth            *auth.Service
	orderService    *order.Service
	withdrawService *withdraw.Service
	balanceService  *balance.Service
}

// NewHandler - constructor
func NewHandler(
	logger logger.Interface,
	cfg config.Config,
	auth *auth.Service,
	orderService *order.Service,
	withdrawService *withdraw.Service,
	balanceService *balance.Service,
) *Handler {
	return &Handler{
		logger:          logger,
		config:          cfg,
		auth:            auth,
		orderService:    orderService,
		withdrawService: withdrawService,
		balanceService:  balanceService,
	}
}
