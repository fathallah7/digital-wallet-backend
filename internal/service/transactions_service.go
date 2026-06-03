package service

import (
	"context"
	"fmt"

	"github.com/fathallah7/wallet-service/internal/apperrors"
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

func (s *TransactionsService) Transfer(ctx context.Context, userID string, req *dto.TransferRequest) error {
	if !req.Amount.IsPositive() {
		return apperrors.ValidationErrors{{Field: "amount", Message: "amount must be greater than zero"}}
	}

	if req.FromWalletID == req.ToWalletID {
		return apperrors.ValidationErrors{{Field: "to_wallet_id", Message: "cannot transfer to the same wallet"}}
	}

	wallet, err := s.walletStore.GetWalletByID(ctx, req.FromWalletID, userID)
	if err != nil {
		return apperrors.ErrWalletNotFound
	}
	if wallet == nil {
		return apperrors.ErrWalletNotFound
	}

	if err := s.transactionsStore.CreateTransfer(ctx, req.FromWalletID, req.ToWalletID, req.Amount); err != nil {
		return err
	}

	return nil
}

func (s *TransactionsService) Deposit(ctx context.Context, userID string, req *dto.DepositRequest) error {
	if !req.Amount.IsPositive() {
		return apperrors.ValidationErrors{{Field: "amount", Message: "amount must be greater than zero"}}
	}

	wallet, err := s.walletStore.GetWalletByID(ctx, req.WalletID, userID)
	if err != nil {
		return apperrors.ErrWalletNotFound
	}
	if wallet == nil {
		return apperrors.ErrWalletNotFound
	}

	if err := s.transactionsStore.Deposit(ctx, req.WalletID, req.Amount); err != nil {
		return fmt.Errorf("deposit: %w", err)
	}

	return nil
}

func (s *TransactionsService) GetUserTransactions(ctx context.Context, userID string) ([]*dto.TransactionResponse, error) {
	transactions, err := s.transactionsStore.GetUserTransactions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get transactions: %w", err)
	}

	var res []*dto.TransactionResponse
	for _, t := range transactions {
		res = append(res, &dto.TransactionResponse{
			ID:           t.ID,
			FromWalletID: t.FromWalletID,
			ToWalletID:   t.ToWalletID,
			Amount:       t.Amount,
			Type:         t.Type,
			Status:       t.Status,
			CreatedAt:    t.CreatedAt,
		})
	}
	return res, nil
}
