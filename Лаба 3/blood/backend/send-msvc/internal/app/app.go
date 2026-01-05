package app

import (
	"context"
	"log"
	"strings"

	"send/internal/core"
	"send/internal/echo/handler"
	"send/internal/echo/server"
	"send/internal/env"
	"send/internal/kafka"
	"send/internal/service"
)

type App struct {
	e        *env.Env
	s        *server.Server
	producer core.Producer
}

func New() (*App, error) {
	a := &App{}

	e, err := env.LoadEnv()
	if err != nil {
		return nil, err
	}
	a.e = e

	brokers := strings.Split(e.KafkaBrokers, ",")
	producer := kafka.NewKafkaGoProducer(brokers, e.KafkaTopic)
	a.producer = producer

	loadCSVService := service.NewLoadCSVService(producer)
	loadCSVHandler := handler.NewLoadCSVHandler(loadCSVService)

	s, err := server.New(loadCSVHandler, e)
	if err != nil {
		_ = producer.Close()
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
	if a.producer != nil {
		_ = a.producer.Close()
	}
	return a.s.Stop(ctx)
}
