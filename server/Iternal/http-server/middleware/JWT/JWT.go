package JWT

import (
	"context"
	"log/slog"
	"net/http"
	"server/Iternal/lib/api/jwt"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		const op string = "jwt-middleware"
		log = log.With(
			slog.String("op", op),
		)

		log.Info("middlewareJWT enabled")

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := jwt.ExtractJWTFromHeader(r)
			if err != nil {
				log.Error("failed extract access token: %v", r)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			claims, err := jwt.ValidateJWT(tokenString)
			if err != nil {
				log.Error("invalid token: %v", err.Error())
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", claims.UserId)
			ctx = context.WithValue(r.Context(), "role", claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
