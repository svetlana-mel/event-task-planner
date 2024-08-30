package middleware

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"log/slog"
	"net/http"

	lib_jwt "github.com/svetlana-mel/event-task-planner/internal/lib/jwt"
	sl "github.com/svetlana-mel/event-task-planner/internal/lib/slog"
	"github.com/svetlana-mel/event-task-planner/internal/server"
)

const TokenHeader = "Authorization"

const EmptyToken = ""

func New(publicKey *ecdsa.PublicKey, log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get(TokenHeader)
			if tokenString == EmptyToken {
				// TODO: add redirect to the /login url or redirect to the main page with login button
				log.Error("no token in request")
				http.Error(w, "Authorization token is required", http.StatusUnauthorized)
				return
			}

			userInfo, err := lib_jwt.ValidateToken(tokenString, publicKey)
			if err != nil {
				if errors.Is(err, lib_jwt.ErrTokenExpired) {
					log.Info("auth middleware: token expired")
					http.Error(w, "Token expired", http.StatusUnauthorized)
					return
				}
				log.Error("invalid token", sl.AddErrorAtribute(err))
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), server.UserContextKey("user"), userInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
