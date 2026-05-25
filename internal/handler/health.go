package handler

import (
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, nil, "running")
}
