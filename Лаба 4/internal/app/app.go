package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lab4/internal/core"
	"lab4/internal/parser"
	"lab4/internal/reader"
	"lab4/internal/saver"
	"lab4/pkg/cfg"
)

type App struct {
	client *http.Client
}

func New() *App {
	return &App{
		client: &http.Client{Timeout: cfg.HTTPTimeout},
	}
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	urls, err := reader.ReadURLsFromExcel()
	if err != nil {
		return err
	}

	results := make([]core.Ship, 0, len(urls))
	for _, url := range urls {
		select {
		case <-ctx.Done():
			fmt.Println("Получен сигнал остановки. Сохраняю уже собранные результаты...")
			goto SAVE
		default:
		}

		fmt.Printf("Обработка: %s\n", url)
		vessel, err := parser.ProcessVesselLink(ctx, a.client, url)
		if err != nil {
			vessel = core.Ship{
				URL:   url,
				Error: err.Error(),
			}
		}
		results = append(results, vessel)
		time.Sleep(cfg.RequestDelay)
	}

SAVE:

	if err := saver.SaveResultsToExcel(results); err != nil {
		return err
	}
	if err := saver.SaveResultsToCSV(results); err != nil {
		return err
	}

	fmt.Println("Обработка завершена. Результаты сохранены в result.xlsx и result.csv")
	return nil
}
