package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransferRequest struct {
	FromWalletID string          `json:"from_wallet_id"`
	ToWalletID   string          `json:"to_wallet_id"`
	Amount       decimal.Decimal `json:"amount"`
}

type DepositRequest struct {
	WalletID string          `json:"wallet_id"`
	Amount   decimal.Decimal `json:"amount"`
}

type TransactionResponse struct {
	ID           string          `json:"id"`
	FromWalletID *string         `json:"from_wallet_id"`
	ToWalletID   *string         `json:"to_wallet_id"`
	Amount       decimal.Decimal `json:"amount"`
	Type         string          `json:"type"`
	Status       string          `json:"status"`
	CreatedAt    time.Time       `json:"created_at"`
}
