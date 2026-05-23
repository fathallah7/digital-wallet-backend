package router

import (
	"net/http"

	"github.com/fathallah7/wallet-service/internal/handler"
)

func Setup() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.HealthHandler)

	return mux
}
