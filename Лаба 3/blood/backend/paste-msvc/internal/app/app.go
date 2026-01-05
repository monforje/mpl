package app

import (
	"context"
	"log"
	"paste/internal/core"
	"paste/internal/echo/handler"
	"paste/internal/echo/server"
	"paste/internal/env"
	"paste/internal/kafka"
	"paste/internal/postgres"
	"paste/internal/service"
	"strings"
)

type App struct {
	e          *env.Env
	pg         *postgres.Postgres
	consumer   core.Consumer
	server     *server.Server
	cancelFunc context.CancelFunc
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

	scanRepo := postgres.NewScanRepository(pg)
	scanService := service.NewScanService(scanRepo)

	brokers := strings.Split(e.KafkaBrokers, ",")
	consumer := kafka.NewKafkaGoConsumer(brokers, e.KafkaTopic, "paste", scanService)
	a.consumer = consumer

	scanHandler := handler.NewScanHandler(scanService)
	srv, err := server.New(scanHandler, e)
	if err != nil {
		_ = consumer.Close()
		pg.Close()
		return nil, err
	}
	a.server = srv

	log.Println("Приложение инициализировано")
	return a, nil
}

func (a *App) Run() error {
	log.Println("Приложение запущено")

	ctx, cancel := context.WithCancel(context.Background())
	a.cancelFunc = cancel

	go func() {
		if kafkaConsumer, ok := a.consumer.(interface{ Start(context.Context) error }); ok {
			log.Println("Запуск Kafka consumer...")
			if err := kafkaConsumer.Start(ctx); err != nil && err != context.Canceled {
				log.Printf("Kafka consumer stopped with error: %v", err)
			}
		}
	}()

	return a.server.Start(a.e.Port)
}

func (a *App) Stop(ctx context.Context) error {
	if a.cancelFunc != nil {
		a.cancelFunc()
	}
	if a.server != nil {
		_ = a.server.Stop(ctx)
	}
	if a.consumer != nil {
		_ = a.consumer.Close()
	}
	if a.pg != nil {
		a.pg.Close()
	}
	return nil
}
