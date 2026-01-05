package main

import (
	"auth/internal/app"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		if err := a.Run(); err != nil {
			log.Printf("server stopped with error: %v", err)
		}
	}()

	log.Println("Приложение запущено")

	<-ctx.Done()
	log.Println("Получен shutdown сигнал")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.Stop(shutdownCtx); err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}

	log.Println("Приложение остановлено")
}
