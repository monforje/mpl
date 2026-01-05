package router

import (
	"send/internal/echo/handler"
	"send/internal/echo/middleware"

	"github.com/labstack/echo/v4"
)

func Setup(
	s *echo.Echo,
	loadCSVHandler *handler.LoadCSVHandler,
	jwtSecret string,
) {
	s.POST("/upload", loadCSVHandler.LoadCSV, middleware.JWTMiddleware(jwtSecret))
}
