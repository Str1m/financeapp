package users

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func pswdVerifedTest(pswd string) bool {
	if len(pswd) < 5 {
		return false
	}
	return true
}

func getHash(pswd string) string {
	h := sha256.New()
	h.Write([]byte(pswd))
	bs := fmt.Sprintf("%x", h.Sum(nil))
	return bs
}

func (s *Service) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// TODO Сделать корректную проверку пароля
	if !pswdVerifedTest(user.Password) {
		http.Error(w, "Password must be more than 5 symbols", http.StatusBadRequest)
		return
	}

	pswdHash := getHash(user.Password)
	user.Password = pswdHash

	user.CreatedAt = time.Now()
	if err := s.SaveUserToDB(&user); err != nil {
		log.Printf("Failed to insert user: %v", err)
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (s *Service) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idStr = strings.TrimSpace(idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid user id: %v", err)
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	user, err := s.GetUserWithID(id)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (s *Service) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		log.Printf("Invalid user id: %v", err)
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	delete(updates, "id")

	if len(updates) == 0 {
		http.Error(w, "Nothing to update", http.StatusBadRequest)
		return
	}

	if err := s.UpdateUserField(id, updates); err != nil {
		log.Printf("Error updating user fields: %v", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}
