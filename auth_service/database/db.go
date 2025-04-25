package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Init() {
	// Используем отдельную переменную для базы пользователей
	dbURL := os.Getenv("AUTH_DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ AUTH_DATABASE_URL не задан в .env")
	}

	var err error
	DB, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к базе данных: %v", err)
	}

	log.Println("✅ Подключение к базе данных AUTH успешно")
}

