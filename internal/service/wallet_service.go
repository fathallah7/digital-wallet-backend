package service

import (
	"context"
	"strings"
	"time"

	"github.com/fathallah7/wallet-service/internal/dto"
	"github.com/fathallah7/wallet-service/internal/store"
)

type WalletService struct {
	walletStore *store.WalletStore
}

func NewWalletService(walletStore *store.WalletStore) *WalletService {
	return &WalletService{
		walletStore: walletStore,
	}
}

func (s *WalletService) CreateWallet(ctx context.Context, req *dto.WalletRequest) (*dto.WalletResponse, map[string]string) {

	if err := validateCreateWalletRequest(req); len(err) > 0 {
		return nil, err
	}

	walletCount, err := s.walletStore.GetUserWalletCount(ctx, req.UserID)
	if err != nil {
		return nil, map[string]string{"wallet_count": "failed to get wallet count"}
	}
	if walletCount >= 3 {
		return nil, map[string]string{"wallet_count": "maximum number of wallets reached (3)"}
	}

	walletId, err := s.walletStore.CreateWallet(ctx, req)
	if err != nil {
		return nil, map[string]string{"create_wallet": "failed to create wallet"}
	}

	return &dto.WalletResponse{
		ID:        walletId,
		UserID:    req.UserID,
		Name:      req.Name,
		Balance:   0,
		IsDefault: false,
		CreatedAt: time.Now(),
	}, nil
}

func (s *WalletService) GetUserWallets(ctx context.Context, userID string) ([]*dto.WalletResponse, map[string]string) {
	if strings.TrimSpace(userID) == "" {
		return nil, map[string]string{"user_id": "user_id is required"}
	}

	wallets, err := s.walletStore.GetUserWallets(ctx, userID)
	if err != nil {
		return nil, map[string]string{"get_wallets": "failed to get wallets"}
	}

	return wallets, nil
}

func validateCreateWalletRequest(req *dto.WalletRequest) map[string]string {
	errors := make(map[string]string)

	if strings.TrimSpace(req.Name) == "" {
		errors["name"] = "name is required"
	}
	if len(strings.TrimSpace(req.Name)) < 3 {
		errors["name"] = "name must be at least 3 characters"
	}
	if strings.TrimSpace(req.UserID) == "" {
		errors["user_id"] = "user_id is required"
	}

	return errors
}
