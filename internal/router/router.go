package router

import (
	"net/http"

	"github.com/fathallah7/wallet-service/internal/handler"
	"github.com/fathallah7/wallet-service/internal/middleware"
)

func Setup(h *handler.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	//
	mux.HandleFunc("GET /health", h.HealthHandler)

	// Auth routes
	mux.HandleFunc("POST /auth/register", h.Register)
	mux.HandleFunc("POST /auth/login", h.Login)

	// Wallet routes
	mux.Handle("POST /wallet", middleware.AuthMiddleware(http.HandlerFunc(h.CreateWallet)))
	mux.Handle("GET /wallets", middleware.AuthMiddleware(http.HandlerFunc(h.GetUserWallets)))
	mux.Handle("GET /wallet/{wallet_id}", middleware.AuthMiddleware(http.HandlerFunc(h.GetWalletByID)))
	mux.Handle("PUT /wallet/{wallet_id}/default", middleware.AuthMiddleware(http.HandlerFunc(h.SetDefaultWallet)))

	// Transaction routes
	mux.Handle("POST /transactions/transfer", middleware.AuthMiddleware(http.HandlerFunc(h.Transfer)))
	mux.Handle("POST /transactions/deposit", middleware.AuthMiddleware(http.HandlerFunc(h.Deposit)))
	return mux
}
