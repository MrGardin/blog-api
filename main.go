package main

import (
	"blog-api/internal/storage/postgres"
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	connString := "postgres://postgres:admin@localhost:5432/blog-api?sslmode=disable"

	config := postgres.Config{
		URL:             connString,
		MaxConns:        20,
		MinConns:        5,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: time.Minute * 30,
	}

	pool, err := postgres.New(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create new pool")
	}
	defer pool.Close()

	//Тестовый запрос версии
	var version string
	err = pool.QueryRow(context.Background(), "SELECT version();").Scan(&version)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	fmt.Println(version)

	// Смотрим статистику пула (для наглядности)
	stats := pool.Stat()
	fmt.Printf("📈 Pool stats: TotalConns(%d) AcquiredConns(%d)\n", stats.TotalConns(), stats.AcquiredConns())
}
