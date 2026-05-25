package handler

import (
	"database/sql"

	"github.com/fathallah7/wallet-service/internal/store"
)

type Handler struct {
	db        *sql.DB
	userStore *store.UserStore
}

func New(db *sql.DB) *Handler {
	return &Handler{
		db:        db,
		userStore: store.NewUserStore(db),
	}
}
