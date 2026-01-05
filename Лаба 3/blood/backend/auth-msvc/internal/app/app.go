package app

import (
	"auth/internal/echo/handler"
	"auth/internal/echo/server"
	"auth/internal/env"
	"auth/internal/postgres"
	"auth/internal/service"
	"context"
	"log"
)

type App struct {
	e           *env.Env
	pg          *postgres.Postgres
	s           *server.Server
	authService *service.AuthService
}

func New() (*App, error) {
	a := &App{}

	e, err := env.LoadEnv()
	if err != nil {
		return nil, err
	}
	a.e = e

	ctx := context.Background()
	pg, err := postgres.New(ctx, e.PostgresURI)
	if err != nil {
		return nil, err
	}
	a.pg = pg

	userRepo := postgres.NewUserRepository(pg)
	authService := service.NewAuthService(userRepo, e.JWTSecret)
	a.authService = authService

	authHandler := handler.NewAuthHandler(authService)

	s, err := server.New(authHandler, authService)
	if err != nil {
		pg.Close()
		return nil, err
	}
	a.s = s

	log.Println("Приложение инициализировано")
	return a, nil
}

func (a *App) Run() error {
	return a.s.Start(a.e.Port)
}

func (a *App) Stop(ctx context.Context) error {
	if a.pg != nil {
		a.pg.Close()
	}
	return a.s.Stop(ctx)
}
