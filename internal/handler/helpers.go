package handler

import (
	"encoding/json"
	"log"
	"net/http"
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

// success
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

// error
func WriteError(w http.ResponseWriter, status int, errors any, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
