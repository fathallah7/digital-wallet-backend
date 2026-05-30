package handler

import (
	"database/sql"

	"github.com/fathallah7/wallet-service/internal/service"
	"github.com/fathallah7/wallet-service/internal/store"
)

type Handler struct {
	db                  *sql.DB
	authService         *service.AuthService
	walletService       *service.WalletService
	transactionsService *service.TransactionsService
}

func New(db *sql.DB) *Handler {
	userStore := store.NewUserStore(db)
	walletStore := store.NewWalletStore(db)
	transactionStore := store.NewTransactionsStore(db)

	authService := service.NewAuthService(userStore)
	walletService := service.NewWalletService(walletStore)
	transactionsService := service.NewTransactionsService(transactionStore, walletStore)

	return &Handler{
		db:                  db,
		authService:         authService,
		walletService:       walletService,
		transactionsService: transactionsService,
	}
}
