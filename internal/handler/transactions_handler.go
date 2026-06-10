package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/fathallah7/wallet-service/internal/dto"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v86"
	"github.com/stripe/stripe-go/v86/webhook"
)

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	if userID == "" {
		WriteError(w, http.StatusBadRequest, nil, "unauthorized")
		return
	}

	var req dto.TransferRequest
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, nil, "invalid request body")
		return
	}

	if err := h.transactionsService.Transfer(r.Context(), userID, &req); WriteServiceError(w, err) {
		return
	}

	WriteJSON(w, http.StatusOK, nil, "transfer successful")
}

func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	if userID == "" {
		WriteError(w, http.StatusBadRequest, nil, "unauthorized")
		return
	}

	var req dto.DepositRequest
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, nil, "invalid request body")
		return
	}

	checkoutURL, err := h.transactionsService.Deposit(r.Context(), userID, &req)
	if WriteServiceError(w, err) {
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"checkout_url": checkoutURL}, "deposit successful")
}

func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	if userID == "" {
		WriteError(w, http.StatusBadRequest, nil, "unauthorized")
		return
	}

	res, err := h.transactionsService.GetUserTransactions(r.Context(), userID)
	if WriteServiceError(w, err) {
		return
	}

	WriteJSON(w, http.StatusOK, res, "transactions retrieved successfully")
}

func (h *Handler) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	signature := r.Header.Get("Stripe-Signature")

	event, err := webhook.ConstructEventWithOptions(
		payload,
		signature,
		endpointSecret,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)

	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		walletID := session.Metadata["wallet_id"]
		amountStr := session.Metadata["amount"]

		amount, err := decimal.NewFromString(amountStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = h.transactionsService.ProcessWebhookDeposit(r.Context(), walletID, amount)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
