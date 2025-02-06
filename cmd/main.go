package main

import (
	db2 "financeapp/internal/db"
	"financeapp/pkg/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := db2.NewDB()
	if err != nil {
		log.Fatalf("Error creating database connection: %v", err)
	}
	defer db.Close()

	userService := users.NewService(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/users", userService.CreateUser)
	http.ListenAndServe(":8080", r)
}
