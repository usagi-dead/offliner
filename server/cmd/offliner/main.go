package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	jwt2 "server/api/lib/jwt"
	"server/internal/s3Storage"

	"github.com/go-chi/httprate"
	swag "github.com/swaggo/http-swagger/v2"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"server/api/lib/emailsender"
	middleJWT "server/api/middleware/jwt"
	middlelog "server/api/middleware/logger"
	_ "server/docs"
	"server/internal/cache"
	"server/internal/config"
	"server/internal/features/user/data"
	"server/internal/features/user/handlers"
	"server/internal/features/user/services"
	"server/internal/storage"
	"syscall"
	"time"
)

type App struct {
	Storage     *storage.Storage
	Cache       *cache.Cache
	S3          *s3Storage.S3Storage
	EmailSender *emailsender.EmailSender
	Router      chi.Router
	Log         *slog.Logger
	Cfg         *config.Config
}

func NewApp(cfg *config.Config, log *slog.Logger) (*App, error) {

	Storage, err := storage.NewStorage(cfg.DbConfig)
	if err != nil {
		return nil, fmt.Errorf("db init failed: %w", err)
	}

	if err := Storage.ApplyMigrations(cfg.DbConfig); err != nil {
		return nil, fmt.Errorf("apply migrations failed: %w", err)
	}

	Cache, err := cache.NewCache(cfg.CacheConfig)
	if err != nil {
		return nil, fmt.Errorf("cache init failed: %w", err)
	}

	s3, err := s3Storage.NewS3Storage(cfg.S3Config)
	if err != nil {
		return nil, fmt.Errorf("s3 init failed: %w", err)
	}

	eSender, err := emailsender.New(cfg.SMTPConfig)
	if err != nil {
		return nil, fmt.Errorf("email sender init failed: %w", err)
	}

	router := chi.NewRouter()
	return &App{Storage: Storage, Cache: Cache, S3: s3, EmailSender: eSender, Router: router, Log: log, Cfg: cfg}, nil
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
	app.Log.Info("docs " + "http://localhost:8080/swagger/index.html#/")

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

	app.Router.Use(
		middleware.Recoverer,
		middleware.RequestID,
		middlelog.New(app.Log),
		middleware.URLFormat,
	)

	var jwtMiddleware = middleJWT.New(app.Log, jwt2.NewJWTHandler(&app.Cfg.JWTConfig))

	//Swagger UI endpoint
	app.Router.Get("/swagger/*", swag.Handler(
		swag.URL("http://localhost:8080/swagger/doc.json"),
	))

	apiVersion := "/v1"

	// Группа для аунтификации
	UserDB := data.NewUserDB(app.Storage.Db)
	UserCache := data.NewUserCache(app.Cache)
	UserS3 := data.NewUserS3(app.S3)
	UserData := data.NewUserQuery(UserDB, UserCache, UserS3)
	UserService := services.NewUserUseCase(UserData, jwt2.NewJWTHandler(&app.Cfg.JWTConfig), app.Log, app.EmailSender)
	UserHandler := handlers.NewUserClient(app.Log, UserService)

	app.Router.Route(apiVersion+"/auth", func(r chi.Router) {
		r.Post("/sign-up", UserHandler.SignUp)
		r.Post("/sign-in", UserHandler.SignIn)
		r.Post("/refresh-token", UserHandler.RefreshToken)
		r.Get("/{provider}", UserHandler.Oauth)
		r.Get("/{provider}/callback", UserHandler.OauthCallback)
		r.Post("/logout", UserHandler.LogoutHandler)
	})

	app.Router.Route(apiVersion+"/confirm", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(httprate.Limit(1, 1*time.Minute, httprate.WithKeyFuncs(httprate.KeyByIP)))
			r.Post("/send-email-code", UserHandler.SendConfirmedEmailCode)
		})
		r.Put("/email", UserHandler.EmailConfirmed)
	})

	// Группа для пользовательских маршрутов (требует авторизации)
	app.Router.Route(apiVersion+"/user", func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Put("/edit", UserHandler.UpdateUserHandler)
		r.Get("/me", UserHandler.GetUserHandler)
		r.Delete("/delete", UserHandler.DeleteUserHandler)
	})

	//// Группа для административных маршрутов
	//app.Router.Route("/admin", func(r chi.Router) {
	//	r.Use(middleJWT.NewCache(app.Log))
	//	//r.Use(middleAdmin)
	//})
	//
	////Группа маршутов для супперадминов для создание админов
	//app.Router.Group(func(r chi.Router) {
	//	r.Use(middleJWT.NewCache(app.Log))
	//	//r.Use(WhiteIpList(WhiteList)
	//	//r.Use(middleSuperAdmin)
	//	//r.Post("/admin/create", auth.CrateAdminHandler(app.Storage, app.Log))
	//})
}

// @title Offliner API
// @version 1.0.0
// @description API for online catalog of PC components.
// @contact.name Evdokimov Igor
// @contact.url https://t.me/epelptic
// @BasePath /v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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
