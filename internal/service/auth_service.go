package service

import (
	"context"
	"net/mail"
	"strings"

	"github.com/fathallah7/wallet-service/internal/dto"
	"github.com/fathallah7/wallet-service/internal/model"
	"github.com/fathallah7/wallet-service/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userStore *store.UserStore
}

func NewAuthService(userStore *store.UserStore) *AuthService {
	return &AuthService{userStore: userStore}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, map[string]string) {
	// 1. Validate
	if err := validateRegister(req); len(err) > 0 {
		return nil, err
	}

	// 2. Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, map[string]string{"general": "something went wrong"}
	}
	hashStr := string(hash)

	// 3. Create user
	u := &model.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        &req.Phone,
		PasswordHash: &hashStr,
	}

	if user, err := s.userStore.GetUserByEmail(ctx, u.Email); user != nil || err != nil {
		return nil, map[string]string{"email": "email is already taken"}
	}

	if user, err := s.userStore.GetUserByPhone(ctx, *u.Phone); user != nil || err != nil {
		return nil, map[string]string{"phone": "phone number is already taken"}
	}

	if err := s.userStore.CreateUser(ctx, u); err != nil {
		return nil, map[string]string{"general": "something went wrong"}
	}

	// 4. Generate token
	token, err := generateToken(u.ID)
	if err != nil {
		return nil, map[string]string{"general": "something went wrong"}
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Phone:     u.Phone,
		},
	}, nil
}

// Login
func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, map[string]string) {

	u, err := s.userStore.GetUserByEmail(ctx, req.Email)
	if u == nil || err != nil {
		return nil, map[string]string{"email": "no account found by this email"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*u.PasswordHash), []byte(req.Password)); err != nil {
		return nil, map[string]string{"password": "wrong password"}
	}

	token, err := generateToken(u.ID)
	if err != nil {
		return nil, map[string]string{"general": "something went wrong"}
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Phone:     u.Phone,
		},
	}, nil
}

func validateRegister(req *dto.RegisterRequest) map[string]string {
	errors := make(map[string]string)

	if len(strings.TrimSpace(req.FirstName)) < 2 {
		errors["firstName"] = "first name must be at least 2 characters"
	}
	if len(strings.TrimSpace(req.LastName)) < 2 {
		errors["lastName"] = "last name must be at least 2 characters"
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		errors["email"] = "invalid email format"
	}
	if len(req.Phone) < 10 {
		errors["phone"] = "invalid phone number"
	}
	if len(req.Password) < 8 {
		errors["password"] = "password must be at least 8 characters"
	}

	return errors
}
