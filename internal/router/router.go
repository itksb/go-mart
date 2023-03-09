package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/itksb/go-mart/internal/handler"
	"github.com/itksb/go-mart/internal/middleware"
	"github.com/itksb/go-mart/internal/service/auth"
	"github.com/itksb/go-mart/pkg/logger"
	"net/http"
)

func NewRouter(
	h *handler.Handler,
	l logger.Interface,
	authServ *auth.Service,
) (http.Handler, error) {
	r := chi.NewRouter()
	r.Use(middleware.GzipUnpackMiddleware)
	// auth middleware here...
	r.Use(middleware.GzipMiddleware)

	corsHandler := middleware.NewCors()
	authMiddleware := middleware.NewAuthMiddleware(authServ, l)

	r.Route("/api/user", func(r2 chi.Router) {
		// apply CORS middleware for api routes
		r2.Use(corsHandler)
		r2.Post("/register", h.APIAuthRegister)
		r2.Post("/login", h.APIAuthLogin)

		r2.Group(func(r3 chi.Router) {
			r3.Use(authMiddleware)
			r3.Post("/orders", h.APIOrderLoad)
			r3.Get("/orders", h.APIGetOrders)
			r3.Post("/balance/withdraw", h.APIWithdraw)
			r3.Get("/balance/withdrawals", h.APIWithdrawals)
			r3.Get("/withdrawals", h.APIWithdrawals)
			r3.Get("/balance", h.APIGetUserSum)
		})
	})

	r.MethodFunc(http.MethodGet, "/health", h.HealthCheck)
	return r, nil
}
