package router

import (
	"auth/internal/echo/handler"
	"auth/internal/echo/middleware"
	"auth/internal/service"

	"github.com/labstack/echo/v4"
)

func Setup(
	s *echo.Echo,
	authHandler *handler.AuthHandler,
	authService *service.AuthService,
) {
	api := s.Group("/api")

	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.GET("/me", authHandler.GetMe, middleware.JWTMiddleware(authService))
}
