package main

import (
	"log/slog"
	"os"
	"server/Iternal/config"
)

func main() {

	cfg := config.MustLoad()

	log := SetupLogger(cfg.Env)
	log.Info("starting offliner server", slog.String("env", cfg.Env))

	//TODO: db init

	//TODO: router init

	//TODO: run server

}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
