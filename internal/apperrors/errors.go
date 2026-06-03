package apperrors

import "errors"

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []FieldError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "validation failed"
	}
	return ve[0].Field + ": " + ve[0].Message
}

var (
	ErrEmailTaken          = errors.New("email already taken")
	ErrPhoneTaken          = errors.New("phone already taken")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrWalletLimit         = errors.New("maximum number of wallets reached")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrWalletNotFound      = errors.New("wallet not found")
)
