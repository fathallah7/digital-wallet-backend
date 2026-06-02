package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Name      string          `json:"name"`
	Balance   decimal.Decimal `json:"balance"`
	IsDefault bool            `json:"is_default"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
