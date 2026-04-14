package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Failed to open DB:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("❌ DB not reachable:", err)
	}

	log.Println("✅ Connected to RDS PostgreSQL!")
	createTables()
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id         SERIAL PRIMARY KEY,
		name       VARCHAR(100) NOT NULL,
		email      VARCHAR(100) UNIQUE NOT NULL,
		age        INT,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("❌ Failed to create tables:", err)
	}
	log.Println("✅ Tables ready!")
}
