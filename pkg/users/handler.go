package users

import (
	"encoding/json"
	"financeapp/pkg/auth"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) Handler {
	return Handler{service: service}
}

func isValidPassword(password string) bool {
	return len(password) >= 5
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func extractID(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	idStr = strings.TrimSpace(idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"message": Invalid JSON}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// TODO Сделать корректную проверку пароля
	if !isValidPassword(user.Password) {
		http.Error(w, `{"message": "Password must be at least 5 characters"}`, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(user); err != nil {
		http.Error(w, `{"message": "Invalid input: `+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	hashedPass, err := hashPassword(user.Password)
	if err != nil {
		log.Printf("Error in Hashing password: %v", err)
		http.Error(w, `{"message": Failed to get password hash}`, http.StatusInternalServerError)
		return
	}
	user.Password = hashedPass
	user.CreatedAt = time.Now()
	if err := h.service.SaveUserToDB(&user); err != nil {
		log.Printf("Failed to insert user: %v", err)
		http.Error(w, `{"message": "Failed to insert user"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetUserWithID(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r)
	if err != nil {
		http.Error(w, `{"message": "Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserWithID(id)
	if err != nil {
		http.Error(w, `{"message": "User not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r)
	if err != nil {
		http.Error(w, `{"message": "Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(auth.UserIDKey).(int)
	if !ok || userID != id {
		http.Error(w, `{"message": "Permission denied"}`, http.StatusForbidden)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, `{"message": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	delete(updates, "id")
	if len(updates) == 0 {
		http.Error(w, `{"message": "Nothing to update"}`, http.StatusBadRequest)
		return
	}

	for field, value := range updates {
		switch field {
		case "email":
			if str, ok := value.(string); !ok || validate.Var(str, "required,email") != nil {
				http.Error(w, `{"message": "Invalid email"}`, http.StatusBadRequest)
				return
			}
		case "password":
			str, ok := value.(string)
			if !ok || validate.Var(str, "required,min=6") != nil {
				http.Error(w, `{"message": "Invalid password"}`, http.StatusBadRequest)
				return
			}
			updates["password"], err = hashPassword(str)
			if err != nil {
				log.Printf("Error in Hashing password: %v", err)
				http.Error(w, `{"message": Failed to get password hash}`, http.StatusInternalServerError)
				return
			}
		default:
			delete(updates, field)
		}
	}

	if err := h.service.UpdateUserField(id, updates); err != nil {
		log.Printf("Error updating user fields: %v", err)
		http.Error(w, `{"message": "Failed to update user"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r)
	if err != nil {
		http.Error(w, `{"message": "Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(auth.UserIDKey).(int)
	if !ok || userID != id {
		http.Error(w, `{"message": "Permission denied"}`, http.StatusForbidden)
		return
	}

	rowAffected, err := h.service.DeleteUserFromDB(id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		http.Error(w, `{"message": "Failed to delete user"}`, http.StatusInternalServerError)
		return
	}
	if rowAffected == 0 {
		http.Error(w, `{"message": "User not found"}`, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) LogUser(w http.ResponseWriter, r *http.Request) {
	var logPass struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&logPass); err != nil {
		http.Error(w, `{"message": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, err := h.service.GetUserByEmail(logPass.Email)
	if err != nil {
		http.Error(w, `{"message": "Invalid username or password"}`, http.StatusInternalServerError)
		return
	}
	err = checkPasswordHash(logPass.Password, user.Password)
	if err != nil {
		http.Error(w, `{"message": "Invalid username or password"}`, http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, `{"message": "Could not generate token"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
