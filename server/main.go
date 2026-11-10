package main

import (
	"context"
	"fmt"
	"kaspersky/internal/db"
	"kaspersky/internal/http"
	"kaspersky/internal/threat"
	"log"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	ctx := context.Background()
	pool, err := db.CreateConnection(ctx)
	if err != nil {
		log.Fatal(err)
	}

	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	repo := threat.NewRepository(pool)
	httpHandler := http.NewHTTPHandlers(repo)
	httpServer := http.NewHTTPServer(httpHandler)

	if err := httpServer.StartServer(); err != nil {
		fmt.Println("не удалось запустить http-сервер: ", err)
	}
}
