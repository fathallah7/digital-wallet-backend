package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fathallah7/wallet-service/internal/dto"
)

func (h *Handler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	if userID == "" {
		WriteError(w, http.StatusBadRequest, nil, "user id is required")
		return
	}

	var req dto.WalletRequest
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, nil, "invalid request body")
		return
	}

	req.UserID = userID
	res, err := h.walletService.CreateWallet(r.Context(), &req)
	if WriteServiceError(w, err) {
		return
	}

	WriteJSON(w, http.StatusOK, res, "wallet created successfully")
}

func (h *Handler) GetUserWallets(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	if userID == "" {
		WriteError(w, http.StatusBadRequest, nil, "user id is required")
		return
	}

	res, err := h.walletService.GetUserWallets(r.Context(), userID)
	if WriteServiceError(w, err) {
		return
	}

	WriteJSON(w, http.StatusOK, res, "wallets retrieved successfully")
}

func (h *Handler) GetWalletByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	if userID == "" {
		WriteError(w, http.StatusBadRequest, nil, "user id is required")
		return
	}

	walletID := r.PathValue("wallet_id")
	if walletID == "" {
		WriteError(w, http.StatusBadRequest, nil, "wallet id is required")
		return
	}

	res, err := h.walletService.GetWalletByID(r.Context(), walletID, userID)
	if WriteServiceError(w, err) {
		return
	}

	WriteJSON(w, http.StatusOK, res, "wallet retrieved successfully")
}

func (h *Handler) SetDefaultWallet(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	if userID == "" {
		WriteError(w, http.StatusBadRequest, nil, "user id is required")
		return
	}

	walletID := r.PathValue("wallet_id")
	if walletID == "" {
		WriteError(w, http.StatusBadRequest, nil, "wallet id is required")
		return
	}

	err := h.walletService.SetDefaultWallet(r.Context(), walletID, userID)
	if WriteServiceError(w, err) {
		return
	}

	WriteJSON(w, http.StatusOK, nil, "default wallet set successfully")
}
