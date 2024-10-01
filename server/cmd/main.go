package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/users", func(c echo.Context) error {
		name := c.FormValue("name")
		return c.String(http.StatusOK, "User created: "+name)
	})

	e.Start(":8080")
}
