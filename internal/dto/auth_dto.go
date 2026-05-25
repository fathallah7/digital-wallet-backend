package dto

import (
	"net/mail"
	"strings"
)

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

func (r *RegisterRequest) Validate() map[string]string {
	errors := make(map[string]string)
	if len(strings.TrimSpace(r.FirstName)) < 2 {
		errors["first_name"] = "first name must be at least 2 characters"
	}

	if len(strings.TrimSpace(r.LastName)) < 2 {
		errors["last_name"] = "last name must be at least 2 characters"
	}

	if len(strings.TrimSpace(r.Email)) == 0 {
		errors["email"] = "Email address is required"
	} else if _, err := mail.ParseAddress(r.Email); err != nil {
		errors["email"] = "invalid email format (e.g., user@gmail.com)"
	}

	if len(strings.TrimSpace(r.Phone)) == 0 {
		errors["phone"] = "Phone number is required"
	} else if len(r.Phone) < 10 {
		errors["phone"] = "Phone number must be at least 10 digits"
	}

	if len(strings.TrimSpace(r.Password)) < 6 {
		errors["password"] = "password must be at least 6 characters"
	}

	if len(r.Password) < 8 {
		errors["password"] = "Password must be at least 8 characters long"
	}

	return errors
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginRequest) Validate() map[string]string {
	errors := make(map[string]string)

	if len(strings.TrimSpace(l.Email)) == 0 {
		errors["email"] = "Email is required"
	}
	if len(l.Password) == 0 {
		errors["password"] = "Password is required"
	}

	return errors
}

// UserResponse represents the response body for user data
type UserResponse struct {
	ID        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Phone     *string `json:"phone"`
}

// AuthResponse represents the response body for authentication
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
