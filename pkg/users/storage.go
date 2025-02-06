package users

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func (s *Service) SaveUserToDB(user *User) error {
	query := `INSERT INTO users (email, password, created_at) VALUES ($1, $2, $3) RETURNING id`
	err := s.DB.QueryRow(query, user.Email, user.Password, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetUserWithID(id int) (*User, error) {
	var user User
	query := `SELECT id, email, password, created_at FROM users WHERE id=$1`
	err := s.DB.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *Service) UpdateUserField(id int, updates map[string]interface{}) error {
	query := `UPDATE users SET `
	i := 1
	values := []interface{}{}
	for key, value := range updates {
		if i > 1 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", key, i)
		values = append(values, value)
		i++
	}
	query += fmt.Sprintf(" WHERE id = $%d", i)
	values = append(values, id)

	_, err := s.DB.Exec(query, values...)
	if err != nil {
		log.Printf("Error updating user fields: %v", err)
		return err
	}
	return nil

}
