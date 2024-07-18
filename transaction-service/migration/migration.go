package main

import (
	"fmt"
	"log"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	fmt.Println("Transaction Service Migration Start --")

	dbURL := "clickhouse://clickhouse:9000?debug=true"
	migrationsPath := "file://files"
	fmt.Println("Connecting to database and loading migration files...")

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	fmt.Println("Running migrations...")

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("Migrations ran successfully")
}
