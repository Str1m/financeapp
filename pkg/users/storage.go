package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

const dbTimeout = 1 * time.Second

func (s *Service) SaveUserToDB(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `INSERT INTO users (email, password, created_at) VALUES ($1, $2, $3) RETURNING id`
	err := s.DB.QueryRowContext(ctx, query, user.Email, user.Password, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		log.Printf("SaveUserToDB error: %v", err)
		return err
	}
	return nil
}

func (s *Service) GetUserWithID(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user User
	query := `SELECT id, email, password, created_at FROM users WHERE id=$1`
	err := s.DB.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("GetUserWithID error: %v", err)
		return nil, err
	}
	return &user, nil
}

func (s *Service) UpdateUserField(id int, updates map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `UPDATE users SET `
	args := []interface{}{}
	i := 1

	for key, value := range updates {
		if i > 1 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", key, i)
		args = append(args, value)
		i++
	}

	query += fmt.Sprintf(" WHERE id = $%d", i)
	args = append(args, id)

	_, err := s.DB.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("UpdateUserField error: %v", err)
		return err
	}
	return nil
}

func (s *Service) DeleteUserFromDB(id int) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM users WHERE id=$1`
	result, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("DeleteUserFromDB error: %v", err)
		return 0, err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("DeleteUserFromDB rows affected error: %v", err)
		return 0, err
	}
	return rowAffected, nil
}

func (s *Service) GetUserByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user User
	query := `SELECT id, email, password, created_at FROM users WHERE email=$1`
	if err := s.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("GetUserByEmail error: %v", err)
		return nil, err
	}
	return &user, nil
}
