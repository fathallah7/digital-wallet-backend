package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// success
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// error
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
