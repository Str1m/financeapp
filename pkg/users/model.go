package users

import (
	"encoding/json"
	"financeapp/internal/db"
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type Service struct {
	DB *db.DB
}

func NewService(db *db.DB) *Service {
	return &Service{DB: db}
}

func pswdVerifedTest(pswd string) bool {
	if len(pswd) < 5 {
		return false
	}
	return true
}

func getHash(pswd string) string {
	return pswd
}

func (s *Service) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// TODO Сделать корректную проверку пароля
	if !pswdVerifedTest(user.Password) {
		http.Error(w, "Password must be more than 5 symbols", http.StatusBadRequest)
		return
	}

	// TODO: Хэширование пароля
	pswdHash := getHash(user.Password)
	user.Password = pswdHash

	user.CreatedAt = time.Now()
	if err = s.SaveUserToDB(&user); err != nil {
		log.Printf("Failed to insert user: %v", err)
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
