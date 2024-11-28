package jwt

import (
	"context"
	"log/slog"
	"net/http"
	ck "server/api/lib/contextKeys"
	"server/api/lib/jwt"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	const op string = "jwt-middleware"
	log = log.With(
		slog.String("op", op),
	)

	log.Info("middlewareJWT enabled")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := jwt.ExtractJWTFromHeader(r)
			if err != nil {
				log.Error("failed extract access token: %v", err.Error())
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			claims, err := jwt.ValidateJWT(tokenString)
			if err != nil {
				log.Error("invalid token: %v", err.Error())
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ck.UserIDKey, claims.UserId)
			ctx = context.WithValue(ctx, ck.RoleKey, claims.Role)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
