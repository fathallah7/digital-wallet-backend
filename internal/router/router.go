package router

import (
	"net/http"

	"github.com/fathallah7/wallet-service/internal/handler"
)

func Setup(h *handler.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.HealthHandler)

	return mux
}
