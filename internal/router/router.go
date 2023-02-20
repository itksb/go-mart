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

	r.Route("/api/user", func(r2 chi.Router) {
		// apply CORS middleware for api routes
		r2.Use(NewCors())

		r2.Post("/register", h.APIAuthRegister)

	})

	r.MethodFunc(http.MethodGet, "/health", h.HealthCheck)
	return r, nil
}
