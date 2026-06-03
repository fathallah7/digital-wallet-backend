package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/fathallah7/wallet-service/internal/apperrors"
)

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

func WriteJSON(w http.ResponseWriter, status int, data any, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
	if err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func WriteError(w http.ResponseWriter, status int, errs any, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Message: message,
		Errors:  errs,
	}); err != nil {
		log.Printf("Error encoding JSON error response: %v", err)
	}
}

func WriteServiceError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	var ve apperrors.ValidationErrors
	if errors.As(err, &ve) {
		WriteError(w, http.StatusBadRequest, ve, "validation failed")
		return true
	}

	switch {
	case errors.Is(err, apperrors.ErrEmailTaken),
		errors.Is(err, apperrors.ErrPhoneTaken):
		WriteError(w, http.StatusConflict, nil, err.Error())
	case errors.Is(err, apperrors.ErrInvalidCredentials):
		WriteError(w, http.StatusUnauthorized, nil, err.Error())
	case errors.Is(err, apperrors.ErrWalletLimit),
		errors.Is(err, apperrors.ErrInsufficientBalance):
		WriteError(w, http.StatusBadRequest, nil, err.Error())
	case errors.Is(err, apperrors.ErrWalletNotFound):
		WriteError(w, http.StatusNotFound, nil, err.Error())
	default:
		log.Printf("Unexpected error: %v", err)
		WriteError(w, http.StatusInternalServerError, nil, "something went wrong")
	}
	return true
}
