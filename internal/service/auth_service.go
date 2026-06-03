package service

import (
	"context"
	"fmt"
	"net/mail"
	"strings"

	"github.com/fathallah7/wallet-service/internal/apperrors"
	"github.com/fathallah7/wallet-service/internal/dto"
	"github.com/fathallah7/wallet-service/internal/model"
	"github.com/fathallah7/wallet-service/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userStore *store.UserStore
	jwtSecret []byte
}

func NewAuthService(userStore *store.UserStore, jwtSecret []byte) *AuthService {
	return &AuthService{userStore: userStore, jwtSecret: jwtSecret}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	if errs := validateRegister(req); len(errs) > 0 {
		return nil, errs
	}

	user, err := s.userStore.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("check email: %w", err)
	}
	if user != nil {
		return nil, apperrors.ErrEmailTaken
	}

	if req.Phone != "" {
		user, err = s.userStore.GetUserByPhone(ctx, req.Phone)
		if err != nil {
			return nil, fmt.Errorf("check phone: %w", err)
		}
		if user != nil {
			return nil, apperrors.ErrPhoneTaken
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	hashStr := string(hash)

	var phone *string
	if req.Phone != "" {
		phone = &req.Phone
	}

	u := &model.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        phone,
		PasswordHash: &hashStr,
	}

	if err := s.userStore.CreateUser(ctx, u); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	token, err := generateToken(u.ID, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
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

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	u, err := s.userStore.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if u == nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*u.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	token, err := generateToken(u.ID, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
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

func validateRegister(req *dto.RegisterRequest) apperrors.ValidationErrors {
	var errors apperrors.ValidationErrors

	if len(strings.TrimSpace(req.FirstName)) < 2 {
		errors = append(errors, apperrors.FieldError{Field: "first_name", Message: "first name must be at least 2 characters"})
	}
	if len(strings.TrimSpace(req.LastName)) < 2 {
		errors = append(errors, apperrors.FieldError{Field: "last_name", Message: "last name must be at least 2 characters"})
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		errors = append(errors, apperrors.FieldError{Field: "email", Message: "invalid email format"})
	}
	if req.Phone != "" && len(req.Phone) < 10 {
		errors = append(errors, apperrors.FieldError{Field: "phone", Message: "invalid phone number"})
	}
	if len(req.Password) < 8 {
		errors = append(errors, apperrors.FieldError{Field: "password", Message: "password must be at least 8 characters"})
	}

	return errors
}
