package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"server/Iternal/cache"
	"server/Iternal/config"
	"server/Iternal/http-server/handlers/auth"
	"server/Iternal/http-server/middleware/emailConfirmed"
	middleJWT "server/Iternal/http-server/middleware/jwt"
	middlelog "server/Iternal/http-server/middleware/logger"
	resp "server/Iternal/lib/api/response"
	"server/Iternal/lib/emailsender"
	"server/Iternal/storage"
	"syscall"
	"time"
)

type App struct {
	Storage     *storage.Storage
	Cache       *cache.Cache
	EmailSender *emailsender.EmailSender
	Router      chi.Router
	Log         *slog.Logger
	Cfg         *config.Config
}

func NewApp(cfg *config.Config, log *slog.Logger) (*App, error) {
	Storage, err := storage.New(cfg.DbConfig)
	if err != nil {
		return nil, fmt.Errorf("db init failed: %w", err)
	}

	Cache, err := cache.New(cfg.CacheConfig)
	if err != nil {
		return nil, fmt.Errorf("cache init failed: %w", err)
	}

	eSender, err := emailsender.New()
	if err != nil {
		return nil, fmt.Errorf("email sender init failed: %w", err)
	}

	router := chi.NewRouter()
	return &App{Storage: Storage, Cache: Cache, EmailSender: eSender, Router: router, Log: log, Cfg: cfg}, nil
}

func (app *App) Start() error {
	srv := &http.Server{
		Addr:         app.Cfg.HttpServerConfig.Address,
		Handler:      app.Router,
		ReadTimeout:  app.Cfg.HttpServerConfig.Timeout,
		WriteTimeout: app.Cfg.HttpServerConfig.Timeout,
		IdleTimeout:  app.Cfg.HttpServerConfig.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Log.Error("server error", slog.String("error", err.Error()))
		}
	}()

	app.Log.Info("server started", slog.String("Addr", app.Cfg.HttpServerConfig.Address))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	app.Log.Info("server stopped gracefully")
	return nil
}

func (app *App) SetupRoutes() {
	app.Router.Use(middleware.RequestID)
	app.Router.Use(middlelog.New(app.Log))
	app.Router.Use(middleware.Recoverer)
	app.Router.Use(middleware.URLFormat)

	// Группа для аунтификации
	app.Router.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", auth.SignUpHandler(app.EmailSender, app.Cache, app.Storage, app.Log))
		r.Post("/sign-in", auth.SignInHandler(app.Storage, app.Log))
		r.Post("/confirmed-email", auth.EmailConfirmedHandler(app.Cache, app.Storage, app.Log))
		//r.Post("/complete-profile", auth.CompleteProfileHandler(app.Storage, app.Log))
		r.Get("/refresh-token", auth.RefreshTokenHandler(app.Storage, app.Log))
		r.Get("/{provider}", auth.OauthHandler(app.Cache, app.Log))
		r.Get("/{provider}/callback", auth.OauthCallbackHandler(app.Cache, app.Log))
	})

	// Группа для пользовательских маршрутов (требует авторизации)
	app.Router.Group(func(r chi.Router) {
		r.Use(middleJWT.New(app.Log))
		r.Use(emailConfirmed.New(app.Storage, app.Log))
		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			render.JSON(w, r, resp.OK())
		})
	})

	// Группа для административных маршрутов
	app.Router.Route("/admin", func(r chi.Router) {
		r.Use(middleJWT.New(app.Log))
		//r.Use(middleAdmin)
	})

	//Группа маршутов для супперадминов для создание админов
	app.Router.Group(func(r chi.Router) {
		r.Use(middleJWT.New(app.Log))
		//r.Use(WhiteIpList(WhiteList)
		//r.Use(middleSuperAdmin)
		//r.Post("/admin/create", auth.CrateAdminHandler(app.Storage, app.Log))
	})
}

func main() {
	cfg := config.MustLoad()
	log := SetupLogger(cfg.Env)

	app, err := NewApp(cfg, log)
	if err != nil {
		log.Error("app init failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	app.SetupRoutes()

	if err := app.Start(); err != nil {
		log.Error("server error", slog.String("error", err.Error()))
		os.Exit(1)
	}
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
