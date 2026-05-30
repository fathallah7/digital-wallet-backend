package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fathallah7/wallet-service/internal/dto"
)

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	if userID == "" {
		WriteError(w, http.StatusBadRequest, nil, "unauthorized")
		return
	}

	var req dto.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, nil, "invalid request body")
		return
	}

	if err := h.transactionsService.Transfer(r.Context(), userID, &req); err != nil {
		WriteError(w, http.StatusBadRequest, err, "transfer failed")
		return
	}

	WriteJSON(w, http.StatusOK, nil, "transfer successful")
}

func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req dto.DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, nil, "invalid request body")
		return
	}

	if err := h.transactionsService.Deposit(r.Context(), userID, &req); err != nil {
		WriteError(w, http.StatusBadRequest, err, "deposit failed")
		return
	}

	WriteJSON(w, http.StatusOK, nil, "deposit successful")
}
