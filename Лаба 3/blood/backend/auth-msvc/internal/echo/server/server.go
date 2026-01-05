package server

import (
	"auth/internal/echo/handler"
	"auth/internal/echo/router"
	"auth/internal/echo/validator"
	"auth/internal/service"
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e *echo.Echo
}

func New(
	authHandler *handler.AuthHandler,
	authService *service.AuthService,
) (*Server, error) {
	s := &Server{}

	s.e = echo.New()
	s.e.HideBanner = true
	s.e.HidePort = true

	s.e.Validator = validator.New()

	s.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://127.0.0.1:5173", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.Setup(
		s.e,
		authHandler,
		authService,
	)

	log.Println("Echo сервер инициализирован")

	return s, nil
}

func (s *Server) Start(port string) error {
	log.Println("Запускаем Echo сервер на порту " + port)
	if err := s.e.Start(":" + port); err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("Останавливаем Echo сервер")
	s.e.Shutdown(ctx)
	return nil
}
