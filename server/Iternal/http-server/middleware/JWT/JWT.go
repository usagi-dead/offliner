package JWT

import (
	"context"
	"log/slog"
	"net/http"
	"server/Iternal/lib/api/jwt"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middlewareJWT"),
		)

		log.Info("middlewareJWT enabled")

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := jwt.ExtractJWTFromHeader(r)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			claims, err := jwt.ValidateJWT(tokenString)
			if err != nil {
				http.Error(w, "invalid token"+err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", claims.UserId)
			ctx = context.WithValue(r.Context(), "role", claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
