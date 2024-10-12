package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"server/Iternal/Storage"
	"server/Iternal/config"
	"server/Iternal/http-server/handlers/auth"
	middleJWT "server/Iternal/http-server/middleware/jwt"
	middlelog "server/Iternal/http-server/middleware/logger"
)

func main() {

	cfg := config.MustLoad()

	log := SetupLogger(cfg.Env)

	storage, err := Storage.New(cfg.DbConfig)
	if err != nil {
		log.Error("db init field: " + err.Error())
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middlelog.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/sign-up", auth.SignUpHandler(storage, log))
	router.Post("/sign-in", auth.SignInHandler(storage, log))
	//router.GET("/refresh-token")

	router.Group(func(r chi.Router) {
		r.Use(middleJWT.New(log))
		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"hello": "world!",
			})
		})
	})

	log.Info("server starting", slog.String("Addr", cfg.HttpServer.Address))

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Error("http server error: " + err.Error())
	}

	log.Info("server stopped")
}

func SetupLogger(env string) (log *slog.Logger) {
	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	return log
}
