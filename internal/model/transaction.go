package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID           string          `json:"id"`
	FromWalletID *string         `json:"from_wallet_id"`
	ToWalletID   *string         `json:"to_wallet_id"`
	Amount       decimal.Decimal `json:"amount"`
	Type         string          `json:"type"`
	Status       string          `json:"status"`
	CreatedAt    time.Time       `json:"created_at"`
}

const (
	TransactionTypeDeposit  = "deposit"
	TransactionTypeTransfer = "transfer"
	TransactionTypeWithdraw = "withdrawal"
)

const (
	TransactionStatusPending   = "pending"
	TransactionStatusCompleted = "completed"
	TransactionStatusFailed    = "failed"
)
