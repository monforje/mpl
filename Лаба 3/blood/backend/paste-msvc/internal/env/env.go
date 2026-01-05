package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Port         string
	PostgresURI  string
	KafkaBrokers string
	KafkaTopic   string
	JWTSecret    string
}

func LoadEnv() (*Env, error) {
	e := &Env{}

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}
	e.Port = port

	uri := os.Getenv("POSTGRES_URI")
	if uri == "" {
		uri = "postgres://dan:dan@localhost:5432/blood?sslmode=disable"
	}
	e.PostgresURI = uri

	e.KafkaBrokers = os.Getenv("KAFKA_BROKERS")
	if e.KafkaBrokers == "" {
		e.KafkaBrokers = "localhost:9092"
	}

	e.KafkaTopic = os.Getenv("KAFKA_TOPIC")
	if e.KafkaTopic == "" {
		e.KafkaTopic = "blood.raw"
	}

	e.JWTSecret = os.Getenv("JWT_SECRET")
	if e.JWTSecret == "" {
		e.JWTSecret = "your-secret-key-here"
	}

	log.Println("Загрузили env переменные")

	return e, nil
}
