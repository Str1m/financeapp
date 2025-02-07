package main

import (
	db2 "financeapp/internal/db"
	"financeapp/pkg/auth"
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
	userHandler := users.NewHandler(userService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/users", func(r chi.Router) {
		r.Post("/register", userHandler.CreateUser)
		r.Post("/login", userHandler.LogUser)
		r.Get("/{id}", userHandler.GetUserWithID)

		r.With(auth.Middleware).Patch("/{id}", userHandler.UpdateUser)
		r.With(auth.Middleware).Delete("/{id}", userHandler.DeleteUser)
	})

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
