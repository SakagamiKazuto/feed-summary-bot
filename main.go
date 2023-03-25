package main

import (
	"feed-summary-bot/controller"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.POST("/bot/handler", controller.HandleBotEvents)

	e.Logger.Fatal(e.Start(":8080"))
}
