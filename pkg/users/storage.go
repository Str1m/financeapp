package users

func (s *Service) SaveUserToDB(user *User) error {
	query := `INSERT INTO users (email, password, created_at) VALUES ($1, $2, $3) RETURNING id`
	err := s.DB.QueryRow(query, user.Email, user.Password, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}
