package store

import (
	"context"
	"database/sql"

	"github.com/fathallah7/wallet-service/internal/dto"
)

type WalletStore struct {
	db *sql.DB
}

func NewWalletStore(db *sql.DB) *WalletStore {
	return &WalletStore{db: db}
}

func (s *WalletStore) CreateWallet(ctx context.Context, req *dto.WalletRequest) (string, error) {
	query := `INSERT INTO wallets (user_id, name) VALUES ($1, $2) RETURNING id`

	var id string
	err := s.db.QueryRowContext(ctx, query, req.UserID, req.Name).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *WalletStore) GetUserWalletCount(ctx context.Context, userID string) (int, error) {
	query := `SELECT COUNT(*) FROM wallets WHERE user_id = $1`
	var count int
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

func (s *WalletStore) GetUserWallets(ctx context.Context, userID string) ([]*dto.WalletResponse, error) {
	query := `SELECT id, user_id, name, balance, is_default, created_at FROM wallets WHERE user_id = $1`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []*dto.WalletResponse
	for rows.Next() {
		var wallet dto.WalletResponse
		if err := rows.Scan(&wallet.ID, &wallet.UserID, &wallet.Name, &wallet.Balance, &wallet.IsDefault, &wallet.CreatedAt); err != nil {
			return nil, err
		}
		wallets = append(wallets, &wallet)
	}
	return wallets, nil
}
