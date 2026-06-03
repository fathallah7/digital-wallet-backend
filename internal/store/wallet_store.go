package store

import (
	"context"
	"database/sql"

	"github.com/fathallah7/wallet-service/internal/model"
)

type WalletStore struct {
	db *sql.DB
}

func NewWalletStore(db *sql.DB) *WalletStore {
	return &WalletStore{db: db}
}

func (s *WalletStore) CreateWallet(ctx context.Context, wallet *model.Wallet) error {
	query := `INSERT INTO wallets (user_id, name) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	return s.db.QueryRowContext(ctx, query, wallet.UserID, wallet.Name).Scan(
		&wallet.ID, &wallet.CreatedAt, &wallet.UpdatedAt,
	)
}

func (s *WalletStore) GetUserWalletCount(ctx context.Context, userID string) (int, error) {
	query := `SELECT COUNT(*) FROM wallets WHERE user_id = $1`
	var count int
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

func (s *WalletStore) GetUserWallets(ctx context.Context, userID string) ([]*model.Wallet, error) {
	query := `SELECT id, user_id, name, balance, is_default, created_at, updated_at FROM wallets WHERE user_id = $1`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []*model.Wallet
	for rows.Next() {
		var wallet model.Wallet
		if err := rows.Scan(&wallet.ID, &wallet.UserID, &wallet.Name, &wallet.Balance, &wallet.IsDefault, &wallet.CreatedAt, &wallet.UpdatedAt); err != nil {
			return nil, err
		}
		wallets = append(wallets, &wallet)
	}
	return wallets, nil
}

func (s *WalletStore) GetWalletByID(ctx context.Context, walletID string, userID string) (*model.Wallet, error) {
	query := `SELECT id, user_id, name, balance, is_default, created_at, updated_at FROM wallets WHERE id = $1 AND user_id = $2`
	var wallet model.Wallet
	err := s.db.QueryRowContext(ctx, query, walletID, userID).Scan(
		&wallet.ID, &wallet.UserID, &wallet.Name, &wallet.Balance, &wallet.IsDefault, &wallet.CreatedAt, &wallet.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (s *WalletStore) SetDefaultWallet(ctx context.Context, userID string, walletID string) error {
	query := `UPDATE wallets SET is_default = CASE WHEN id = $1 THEN true ELSE false END, updated_at = NOW() WHERE user_id = $2`
	_, err := s.db.ExecContext(ctx, query, walletID, userID)
	return err
}
