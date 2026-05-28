package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fathallah7/wallet-service/internal/dto"
)

// create wallet
func (h *Handler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	var req dto.WalletRequest

	userId := r.Context().Value("user_id").(string)
	if userId == "" {
		WriteError(w, http.StatusBadRequest, nil, "user id is required")
		return
	}

	req.UserID = userId

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, nil, "invalid request body")
		return
	}

	res, err := h.walletService.CreateWallet(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err, "validation failed")
		return
	}

	WriteJSON(w, http.StatusOK, res, "wallet created successfully")
}

// get user wallets
func (h *Handler) GetUserWallets(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	if userID == "" {
		WriteError(w, http.StatusBadRequest, nil, "")
		return
	}

	res, err := h.walletService.GetUserWallets(r.Context(), userID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err, "validation failed")
		return
	}

	WriteJSON(w, http.StatusOK, res, "wallet retrieved successfully")
}
