package domain

import "errors"

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInsufficientFunds = errors.New("insufficient pool funds")
	ErrInvalidState      = errors.New("invalid state transition")
)
