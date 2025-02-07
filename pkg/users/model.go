package users

import (
	"financeapp/internal/db"
	"github.com/go-playground/validator/v10"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var validate = validator.New()

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=6"`
	CreatedAt time.Time `json:"created_at"`
}

type Service struct {
	DB *db.DB
}

func NewService(db *db.DB) *Service {
	return &Service{DB: db}
}
