package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"server/Iternal/config"
	storage2 "server/Iternal/storage"
)

func main() {
	e := echo.New()

	cfg := config.MustLoad()

	SetupLogger(cfg.Env, e)
	e.Logger.Info(cfg)

	storage, err := storage2.New(cfg.DbPath)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Info(storage.Db.String())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":8080"))

}

func SetupLogger(env string, e *echo.Echo) {
	switch env {
	case "local":
		e.Logger.SetLevel(1)
		e.Logger.SetOutput(os.Stdout)
	case "dev":
		e.Logger.SetLevel(2)
	}
}
