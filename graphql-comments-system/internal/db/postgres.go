package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func InitPostgres() {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN is not set in environment")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to open postgres connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping postgres: %v", err)
	}

	schemaPath := "internal/db/schema.sql"
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("failed to read schema.sql: %v", err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		log.Fatalf("failed to execute schema.sql: %v", err)
	}

	DB = db
	fmt.Println("Connected to PostgreSQL and schema loaded")
}
