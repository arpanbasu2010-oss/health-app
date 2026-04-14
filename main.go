package main

import (
	"go-api/config"
	"go-api/handler"
	"log"

	repoImpl "go-api/repo/impl"
	serviceImpl "go-api/service/impl"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file, using system env vars")
	}

	// Init DB
	config.InitDB()

	// Wire up layers  (Dependency Injection)
	userRepo := repoImpl.NewUserRepo(config.DB)
	userService := serviceImpl.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Router
	r := mux.NewRouter()

	// Routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
