package transaction

import "time"

type Transaction struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	AccountID   int       `json:"account_id"`
	CategoryID  int       `json:"category_id"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
}
