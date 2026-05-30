package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fathallah7/wallet-service/internal/model"
)

type TransactionsStore struct {
	db *sql.DB
}

func NewTransactionsStore(db *sql.DB) *TransactionsStore {
	return &TransactionsStore{db: db}
}

func (s *TransactionsStore) CreateTransfer(ctx context.Context, fromWalletID, toWalletID string, amount float64) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var balance float64
	err = tx.QueryRowContext(ctx, "SELECT balance FROM wallets WHERE id = $1 FOR UPDATE", fromWalletID).Scan(&balance)
	if err != nil {
		return err
	}

	if balance < amount {
		return errors.New("insufficient balance")
	}

	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance = balance - $1 WHERE id = $2", amount, fromWalletID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance = balance + $1 WHERE id = $2", amount, toWalletID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		"INSERT INTO transactions (from_wallet_id, to_wallet_id, amount, type, status) VALUES ($1, $2, $3, 'transfer', 'completed')",
		fromWalletID, toWalletID, amount,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *TransactionsStore) Deposit(ctx context.Context, walletID string, amount float64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`UPDATE wallets SET balance = balance + $1 WHERE id = $2`,
		amount, walletID,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO transactions (to_wallet_id, amount, type, status)
		 VALUES ($1, $2, 'deposit', 'completed')`,
		walletID, amount,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *TransactionsStore) GetUserTransactions(ctx context.Context, userID string) ([]*model.Transaction, error) {
	query := `
		SELECT t.id, t.from_wallet_id, t.to_wallet_id, t.amount, t.type, t.status, t.created_at
		FROM transactions t
		JOIN wallets w ON w.id = t.to_wallet_id OR w.id = t.from_wallet_id
		WHERE w.user_id = $1
		ORDER BY t.created_at DESC
	`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.ID, &t.FromWalletID, &t.ToWalletID, &t.Amount, &t.Type, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}

	return transactions, nil
}
