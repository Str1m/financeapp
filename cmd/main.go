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
		log.Fatalf("Error loading .env file: %v", err)
	}

	db, err := db2.NewDB()
	if err != nil {
		log.Fatalf("Error creating database connection: %v", err)
	}
	defer db.Close()

	userService := users.NewService(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userService.CreateUser)
		r.Get("/{id}", userService.GetUser)
		r.Patch("/{id}", userService.UpdateUser)
	})

	http.ListenAndServe(":8080", r)
}
