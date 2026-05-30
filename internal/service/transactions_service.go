package service

import (
	"context"

	"github.com/fathallah7/wallet-service/internal/dto"
	"github.com/fathallah7/wallet-service/internal/store"
)

type TransactionsService struct {
	transactionsStore *store.TransactionsStore
	walletStore       *store.WalletStore
}

func NewTransactionsService(transactionsStore *store.TransactionsStore, walletStore *store.WalletStore) *TransactionsService {
	return &TransactionsService{
		transactionsStore: transactionsStore,
		walletStore:       walletStore,
	}
}

func (s *TransactionsService) Transfer(ctx context.Context, userID string, req *dto.TransferRequest) map[string]string {
	if req.Amount <= 0 {
		return map[string]string{"error": "amount must be greater than zero"}
	}

	if req.FromWalletID == req.ToWalletID {
		return map[string]string{"error": "cannot transfer to the same wallet"}
	}

	if wallet, err := s.walletStore.GetWalletByID(ctx, req.FromWalletID, userID); err != nil || wallet == nil {
		return map[string]string{"from_wallet": "wallet not found"}
	}

	if err := s.transactionsStore.CreateTransfer(ctx, req.FromWalletID, req.ToWalletID, req.Amount); err != nil {
		if err.Error() == "insufficient balance" {
			return map[string]string{"balance": "insufficient balance"}
		}
		return map[string]string{"general": "transfer failed"}
	}

	return nil
}

func (s *TransactionsService) Deposit(ctx context.Context, userID string, req *dto.DepositRequest) map[string]string {
	if req.Amount <= 0 {
		return map[string]string{"amount": "amount must be greater than 0"}
	}

	wallet, err := s.walletStore.GetWalletByID(ctx, req.WalletID, userID)
	if err != nil || wallet == nil {
		return map[string]string{"wallet": "wallet not found"}
	}

	if err := s.transactionsStore.Deposit(ctx, req.WalletID, req.Amount); err != nil {
		return map[string]string{"general": "deposit failed"}
	}

	return nil
}
