package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/itksb/go-mart/internal/handler"
	"github.com/itksb/go-mart/pkg/logger"
	"net/http"
)

func NewRouter(
	h *handler.Handler,
	l logger.Interface,
) (http.Handler, error) {
	r := chi.NewRouter()
	r.Use(gzipUnpackMiddleware)
	// auth middleware here...
	r.Use(gzipMiddleware)

	// apply CORS middleware for api routes
	r.Use(NewCors())

	return r, nil
}
