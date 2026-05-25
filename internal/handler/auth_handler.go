package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fathallah7/wallet-service/internal/dto"
	"github.com/fathallah7/wallet-service/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// Register handles user registration
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {

	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, nil, "invalid request body")
		return
	}

	if errors := req.Validate(); len(errors) > 0 {
		writeError(w, http.StatusBadRequest, errors, "validation failed")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, nil, "something went wrong")
		return
	}
	hashStr := string(hash)

	u := &model.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        &req.Phone,
		PasswordHash: &hashStr,
	}

	if err := h.userStore.CreateUser(r.Context(), u); err != nil {
		writeError(w, http.StatusInternalServerError, nil, "could not create user")
		return
	}

	token, err := generateToken(u.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, nil, "could not generate token")
		return
	}

	res := dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Phone:     u.Phone,
		},
	}

	writeJSON(w, http.StatusCreated, res, "account created successfully")
}

// Login handles user login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, nil, "invalid request body")
		return
	}

	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		writeError(w, http.StatusBadRequest, validationErrors, "validation failed")
		return
	}

	u, err := h.userStore.GetUserByEmail(r.Context(), req.Email)

	if err != nil {
		writeError(w, http.StatusInternalServerError, nil, "something went wrong")
		return
	}

	if u == nil {
		writeError(w, http.StatusUnauthorized, nil, "no account found by this email")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*u.PasswordHash), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, nil, "wrong password")
		return
	}

	token, err := generateToken(u.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, nil, "could not generate token")
		return
	}

	res := dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Phone:     u.Phone,
		},
	}

	writeJSON(w, http.StatusOK, res, "login successful")
}
