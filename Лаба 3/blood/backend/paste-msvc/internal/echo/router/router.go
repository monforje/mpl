package router

import (
	"paste/internal/echo/handler"
	"paste/internal/echo/middleware"

	"github.com/labstack/echo/v4"
)

func Setup(
	e *echo.Echo,
	scanHandler *handler.ScanHandler,
	jwtSecret string,
) {
	e.GET("/scans", scanHandler.GetMyScans, middleware.JWTMiddleware(jwtSecret))
}
