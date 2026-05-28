package dto

import "time"

type WalletRequest struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type WalletResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Balance   float64   `json:"balance"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
}
