package server

import (
	"context"
	"log"
	"paste/internal/echo/handler"
	"paste/internal/echo/router"
	"paste/internal/env"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e *echo.Echo
}

func New(
	scanHandler *handler.ScanHandler,
	env *env.Env,
) (*Server, error) {
	s := &Server{}

	s.e = echo.New()
	s.e.HideBanner = true
	s.e.HidePort = true

	s.e.Use(middleware.CORS())

	router.Setup(
		s.e,
		scanHandler,
		env.JWTSecret,
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
	return s.e.Shutdown(ctx)
}
