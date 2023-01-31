package handler

import (
	"github.com/itksb/go-mart/internal/config"
	"github.com/itksb/go-mart/pkg/logger"
)

type Handler struct {
	logger logger.Interface
	config config.Config
}

// NewHandler - constructor
func NewHandler(
	logger logger.Interface,
	cfg config.Config,
) *Handler {
	return &Handler{
		logger: logger,
		config: cfg,
	}
}
