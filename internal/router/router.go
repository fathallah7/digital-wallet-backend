package router

import (
	"net/http"

	"github.com/fathallah7/wallet-service/internal/handler"
	"github.com/fathallah7/wallet-service/internal/middleware"
)

func Setup(h *handler.Handler, jwtSecret []byte) *http.ServeMux {
	mux := http.NewServeMux()
	authMW := middleware.NewAuthMiddleware(jwtSecret)

	mux.HandleFunc("GET /health", h.HealthHandler)

	mux.HandleFunc("POST /auth/register", h.Register)
	mux.HandleFunc("POST /auth/login", h.Login)

	mux.Handle("POST /wallet", authMW.Authenticate(http.HandlerFunc(h.CreateWallet)))
	mux.Handle("GET /wallets", authMW.Authenticate(http.HandlerFunc(h.GetUserWallets)))
	mux.Handle("GET /wallet/{wallet_id}", authMW.Authenticate(http.HandlerFunc(h.GetWalletByID)))
	mux.Handle("PUT /wallet/{wallet_id}/default", authMW.Authenticate(http.HandlerFunc(h.SetDefaultWallet)))

	mux.Handle("POST /transactions/transfer", authMW.Authenticate(http.HandlerFunc(h.Transfer)))
	mux.Handle("POST /transactions/deposit", authMW.Authenticate(http.HandlerFunc(h.Deposit)))
	mux.Handle("GET /transactions", authMW.Authenticate(http.HandlerFunc(h.GetTransactions)))

	mux.HandleFunc("POST /webhook/stripe", h.StripeWebhook)

	return mux
}
