package main

import (
	"errors"
	"fmt"
	"os"

	// Library for migrations
	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"

	// Driver for perfirming migrations
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// Driver for getting migrations from files
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	migrationsTable := os.Getenv("MIGRATIONS_TABLE")

	if dbname == "" || user == "" || password == "" {
		panic("database credentials are required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	// Forming a connection string to PostgreSQL
	postgresURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&x-migrations-table=%s",
		user, password, host, port, dbname, migrationsTable)

	// Creating a Migrator
	m, err := migrate.New(
		"file://"+migrationsPath,
		postgresURL,
	)
	if err != nil {
		panic(err)
	}

	// Using migrations
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}

		panic(err)
	}

	fmt.Println("migrations applied successfully")
}
