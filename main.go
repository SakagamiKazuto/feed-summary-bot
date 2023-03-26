package main

import (
	"feed-summary-bot/controller"
	"feed-summary-bot/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	e := echo.New()

	logger.NewLogger()

	e.Use(middleware.Logger())
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "running this server")
	})
	e.POST("/bot/handler", controller.HandleBotEvents)

	e.Logger.Fatal(e.Start(":8080"))
}
