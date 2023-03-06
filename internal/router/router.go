package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/itksb/go-mart/internal/handler"
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
	r.Use(gzipUnpackMiddleware)
	// auth middleware here...
	r.Use(gzipMiddleware)

	corsHandler := NewCors()
	authMiddleware := NewAuthMiddleware(authServ, l)

	r.Route("/api/user", func(r2 chi.Router) {
		// apply CORS middleware for api routes
		r2.Use(corsHandler)
		r2.Post("/register", h.APIAuthRegister)
		r2.Post("/login", h.APIAuthLogin)

		r2.Group(func(r3 chi.Router) {
			r3.Use(authMiddleware)
			//r3.Post("/orders", orderHandler.PostUserOrder())
			/*
				r3.Get("/orders", orderHandler.GetUserOrders())
				r3.Get("/balance", balanceHandler.GetUserSummary())
				r3.Post("/balance/withdraw", withdrawHandler.PostUserBalanceWithdraw())
				r3.Get("/withdrawals", withdrawHandler.GetUserWithdrawals())
			*/
		})
	})

	r.MethodFunc(http.MethodGet, "/health", h.HealthCheck)
	return r, nil
}
