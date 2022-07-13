package api

import (
	"github.com/labstack/echo/v4"
	"github.com/rrobrms/eip712-go-api/handlers"
)

func MainGroup(e *echo.Echo) {
	// Routes
	e.POST("/", handlers.PostEIP712)
	e.POST("/verify", handlers.PostEIP712verify)
}
