package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fathallah7/wallet-service/internal/dto"
)

// Register handles user registration
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {

	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, nil, "invalid request body")
		return
	}

	res, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err, "validation failed")
		return
	}

	writeJSON(w, http.StatusCreated, res, "account created successfully")
}

// Login handles user login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, nil, "invalid request body")
		return
	}

	res, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err, "validation failed")
		return
	}

	writeJSON(w, http.StatusOK, res, "login successful")
}
