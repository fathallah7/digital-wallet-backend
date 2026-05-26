package handler

import (
	"database/sql"

	"github.com/fathallah7/wallet-service/internal/service"
	"github.com/fathallah7/wallet-service/internal/store"
)

type Handler struct {
	db          *sql.DB
	authService *service.AuthService
}

func New(db *sql.DB) *Handler {
	userStore   := store.NewUserStore(db)
	authService := service.NewAuthService(userStore)

	return &Handler{
		db:          db,
		authService: authService,
	}
}
