package middleware

import (
	"context"
	"github.com/itksb/go-mart/internal/service/auth"
	"github.com/itksb/go-mart/internal/service/auth/token"
	"github.com/itksb/go-mart/pkg/logger"
	"net/http"
	"strings"
)

type ctxUserID string

const ctxUser ctxUserID = "user_id"

// NewAuthMiddleware - setup user context
// see examples: https://bash-shell.net/blog/dependency-injection-golang-http-middleware/
func NewAuthMiddleware(authSrv *auth.Service, l logger.Interface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				splitted := strings.Split(
					r.Header.Get("Authorization"),
					" ")
				if len(splitted) != 2 {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				tokenString := splitted[1]
				claims := token.MartClaims{}
				err := authSrv.ParseWithClaims(tokenString, &claims)

				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					l.Infof("auth error: %s", err)
					return
				}

				ctx := context.WithValue(r.Context(), ctxUser, claims.UserID)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
	}
}
