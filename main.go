package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var db *sql.DB

func initDB() {

	_ = godotenv.Load()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || name == "" {
		log.Fatal("Missing required Database Environment Variables (Check your .env or prod.env file)")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, pass, name,
	)

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Invalid DSN format:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database connection failed (Ping):", err)
	}

	fmt.Println("Successfully connected to the database!")

	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL
	);`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Could not create table:", err)
	}
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/users", usersHandler) // New endpoint

	port := ":8080"
	fmt.Println("Server running on", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Example: Create a user
		var u struct{ Name, Email string }
		json.NewDecoder(r.Body).Decode(&u)
		_, err := db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", u.Name, u.Email)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}

	// Example: List users
	rows, _ := db.Query("SELECT name, email FROM users")
	var users []map[string]string
	for rows.Next() {
		var name, email string
		rows.Scan(&name, &email)
		users = append(users, map[string]string{"name": name, "email": email})
	}
	json.NewEncoder(w).Encode(users)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
