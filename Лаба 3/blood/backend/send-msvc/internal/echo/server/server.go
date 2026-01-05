package server

import (
	"context"
	"log"
	"send/internal/echo/handler"
	"send/internal/echo/router"
	"send/internal/env"

	"github.com/labstack/echo/v4"
	middlewareEcho "github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e *echo.Echo
}

func New(
	loadCSVHandler *handler.LoadCSVHandler,
	e *env.Env,
) (*Server, error) {
	s := &Server{}

	s.e = echo.New()
	s.e.HideBanner = true
	s.e.HidePort = true

	s.e.Use(middlewareEcho.CORS())

	router.Setup(
		s.e,
		loadCSVHandler,
		e.JWTSecret,
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
