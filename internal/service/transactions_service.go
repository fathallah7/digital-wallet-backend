package service

import (
	"context"
	"fmt"
	"os"

	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"

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

func (s *TransactionsService) Deposit(ctx context.Context, userID string, req *dto.DepositRequest) (string, error) {
	if !req.Amount.IsPositive() {
		return "", apperrors.ValidationErrors{{Field: "amount", Message: "amount must be greater than zero"}}
	}

	wallet, err := s.walletStore.GetWalletByID(ctx, req.WalletID, userID)
	if err != nil {
		return "", apperrors.ErrWalletNotFound
	}
	if wallet == nil {
		return "", apperrors.ErrWalletNotFound
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	amountInCents := req.Amount.Mul(decimal.NewFromInt(100)).IntPart()

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String("https://github.com/fathallah7"), // TODO: replace with actual success URL
		CancelURL:          stripe.String("https://github.com/fathallah7"), // TODO: replace with actual cancel URL
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String("usd"),
					UnitAmount: stripe.Int64(amountInCents),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String("Wallet Refill"),
						Description: stripe.String("Deposit to wallet"),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},

		Metadata: map[string]string{
			"wallet_id": req.WalletID,
			"user_id":   userID,
			"amount":    req.Amount.String(),
		},
	}

	stripeSession, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("stripe session error: %w", err)
	}

	return stripeSession.URL, nil
}

func (s *TransactionsService) ProcessWebhookDeposit(ctx context.Context, walletID string, amount decimal.Decimal) error {
	return s.transactionsStore.Deposit(ctx, walletID, amount)
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
