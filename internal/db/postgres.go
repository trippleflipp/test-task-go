package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitDB() *sqlx.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}

	if err := applyMigrations(db); err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	fmt.Println("Database migrations applied.")
	return db
}

func applyMigrations(db *sqlx.DB) error {
	query, err := os.ReadFile("migrations/1_init.up.sql")
	if err != nil {
		return fmt.Errorf("Migrations file not found: %w", err)
	}

	_, err = db.Exec(string(query))
	return err
}
