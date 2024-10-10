package main

import (
	"log/slog"
	"os"
	"server/Iternal/Storage"
	"server/Iternal/config"
)

func main() {

	cfg := config.MustLoad()

	log := SetupLogger(cfg.Env)

	storage, err := Storage.New(cfg.DbPath)
	if err != nil {
		log.Error("db init field: " + err.Error())
		os.Exit(1)
	}

	log.Info("db init success: " + storage.Db.String())
	//TODO: init router

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
