package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки .env: %v", err)
	}

	postgresURI := os.Getenv("POSTGRES_URI")
	if postgresURI == "" {
		log.Fatal("POSTGRES_URI не установлен в .env")
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Ошибка установки диалекта: %v", err)
	}

	db, err := sql.Open("pgx", postgresURI)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Ошибка проверки подключения: %v", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("Ошибка выполнения миграций: %v", err)
	}

	log.Println("✅ Все миграции успешно выполнены")
}
