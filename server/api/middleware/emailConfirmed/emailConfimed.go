package emailConfirmed

import (
	"log/slog"
	"net/http"
	ck "server/api/lib/contextKeys"
)

type GetterEmailStatus interface {
	IsEmailConfirmed(UserId int64) (bool, error)
}

func New(getterEmailStatus GetterEmailStatus, log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		const op string = "email-confirmed-middleware"
		log = log.With(
			slog.String("op", op),
		)

		log.Info("email-confirmed-middleware enabled")

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId, ok := r.Context().Value(ck.UserIDKey).(int64)
			if !ok {
				http.Error(w, "user_id not found or invalid", http.StatusBadRequest)
				return
			}

			EmailConfirmed, err := getterEmailStatus.IsEmailConfirmed(userId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !EmailConfirmed {
				http.Error(w, "email not confirmed", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
