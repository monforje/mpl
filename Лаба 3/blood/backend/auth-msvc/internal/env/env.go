package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Port        string
	PostgresURI string
	JWTSecret   string
}

func LoadEnv() (*Env, error) {
	e := &Env{}

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	e.Port = port

	e.PostgresURI = os.Getenv("POSTGRES_URI")
	if e.PostgresURI == "" {
		e.PostgresURI = "postgres://user:password@localhost:5432/blooddb"
	}

	e.JWTSecret = os.Getenv("JWT_SECRET")
	if e.JWTSecret == "" {
		e.JWTSecret = "your-secret-key-here"
	}

	log.Println("Загрузили env переменные")

	return e, nil
}
