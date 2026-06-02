package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type WalletRequest struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type WalletResponse struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Name      string          `json:"name"`
	Balance   decimal.Decimal `json:"balance"`
	IsDefault bool            `json:"is_default"`
	CreatedAt time.Time       `json:"created_at"`
}
