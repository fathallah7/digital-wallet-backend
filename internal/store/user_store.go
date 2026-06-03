package store

import (
	"context"
	"database/sql"

	"github.com/fathallah7/wallet-service/internal/model"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) CreateUser(ctx context.Context, u *model.User) error {
	query := `
		INSERT INTO users (first_name, last_name, email, phone, password_hash)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return s.db.QueryRowContext(ctx, query,
		u.FirstName, u.LastName, u.Email, u.Phone, u.PasswordHash,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, password_hash, created_at, updated_at
		FROM users WHERE email = $1
	`
	u := &model.User{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.FirstName, &u.LastName,
		&u.Email, &u.Phone, &u.PasswordHash,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (s *UserStore) GetUserByPhone(ctx context.Context, phone string) (*model.User, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, password_hash, created_at, updated_at
		FROM users WHERE phone = $1
	`
	u := &model.User{}
	err := s.db.QueryRowContext(ctx, query, phone).Scan(
		&u.ID, &u.FirstName, &u.LastName,
		&u.Email, &u.Phone, &u.PasswordHash,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}
