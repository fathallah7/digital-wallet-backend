package handler

import (
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, HealthResponse{Status: "running"}, "service is healthy")
}
