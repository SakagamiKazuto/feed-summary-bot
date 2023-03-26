package main

import (
	"feed-summary-bot/controller"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "running this server")
	})
	e.POST("/bot/handler", controller.HandleBotEvents)

	e.Logger.Fatal(e.Start(":8080"))
}
