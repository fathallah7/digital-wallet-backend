package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/fathallah7/wallet-service/internal/apperrors"
	"github.com/fathallah7/wallet-service/internal/dto"
	"github.com/fathallah7/wallet-service/internal/model"
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

func (s *WalletService) CreateWallet(ctx context.Context, req *dto.WalletRequest) (*dto.WalletResponse, error) {
	if errs := validateCreateWalletRequest(req); len(errs) > 0 {
		return nil, errs
	}

	walletCount, err := s.walletStore.GetUserWalletCount(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("get wallet count: %w", err)
	}
	if walletCount >= 3 {
		return nil, apperrors.ErrWalletLimit
	}

	wallet := &model.Wallet{
		UserID:    req.UserID,
		Name:      req.Name,
		Balance:   decimal.NewFromInt(0),
		IsDefault: false,
	}

	if err := s.walletStore.CreateWallet(ctx, wallet); err != nil {
		return nil, fmt.Errorf("create wallet: %w", err)
	}

	return &dto.WalletResponse{
		ID:        wallet.ID,
		UserID:    wallet.UserID,
		Name:      wallet.Name,
		Balance:   wallet.Balance,
		IsDefault: wallet.IsDefault,
		CreatedAt: wallet.CreatedAt,
	}, nil
}

func (s *WalletService) GetUserWallets(ctx context.Context, userID string) ([]*dto.WalletResponse, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, apperrors.ValidationErrors{{Field: "user_id", Message: "user id is required"}}
	}

	wallets, err := s.walletStore.GetUserWallets(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get wallets: %w", err)
	}

	var res []*dto.WalletResponse
	for _, w := range wallets {
		res = append(res, &dto.WalletResponse{
			ID:        w.ID,
			UserID:    w.UserID,
			Name:      w.Name,
			Balance:   w.Balance,
			IsDefault: w.IsDefault,
			CreatedAt: w.CreatedAt,
		})
	}
	return res, nil
}

func (s *WalletService) GetWalletByID(ctx context.Context, walletID string, userID string) (*dto.WalletResponse, error) {
	if strings.TrimSpace(walletID) == "" {
		return nil, apperrors.ValidationErrors{{Field: "wallet_id", Message: "wallet id is required"}}
	}

	wallet, err := s.walletStore.GetWalletByID(ctx, walletID, userID)
	if err != nil {
		return nil, apperrors.ErrWalletNotFound
	}

	return &dto.WalletResponse{
		ID:        wallet.ID,
		UserID:    wallet.UserID,
		Name:      wallet.Name,
		Balance:   wallet.Balance,
		IsDefault: wallet.IsDefault,
		CreatedAt: wallet.CreatedAt,
	}, nil
}

func (s *WalletService) SetDefaultWallet(ctx context.Context, walletID string, userID string) error {
	if err := s.walletStore.SetDefaultWallet(ctx, userID, walletID); err != nil {
		return fmt.Errorf("set default wallet: %w", err)
	}
	return nil
}

func validateCreateWalletRequest(req *dto.WalletRequest) apperrors.ValidationErrors {
	var errors apperrors.ValidationErrors

	if strings.TrimSpace(req.Name) == "" {
		errors = append(errors, apperrors.FieldError{Field: "name", Message: "name is required"})
	} else if len(strings.TrimSpace(req.Name)) < 3 {
		errors = append(errors, apperrors.FieldError{Field: "name", Message: "name must be at least 3 characters"})
	}
	if strings.TrimSpace(req.UserID) == "" {
		errors = append(errors, apperrors.FieldError{Field: "user_id", Message: "user id is required"})
	}

	return errors
}
