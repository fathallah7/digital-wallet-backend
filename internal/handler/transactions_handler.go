package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fathallah7/wallet-service/internal/dto"
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

	if err := h.transactionsService.Deposit(r.Context(), userID, &req); WriteServiceError(w, err) {
		return
	}

	WriteJSON(w, http.StatusOK, nil, "deposit successful")
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
