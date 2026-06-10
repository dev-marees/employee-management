// Package apperror defines transport-agnostic sentinel errors that the service
// layer returns and the handler layer maps to HTTP status codes.
package apperror

import "errors"

var (
	// ErrNotFound indicates a requested resource does not exist.
	ErrNotFound = errors.New("resource not found")
	// ErrConflict indicates a uniqueness or state conflict (e.g. duplicate email).
	ErrConflict = errors.New("resource conflict")
	// ErrInvalidInput indicates semantically invalid input not caught by binding.
	ErrInvalidInput = errors.New("invalid input")
	// ErrUnauthorized indicates failed authentication.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden indicates the caller lacks permission.
	ErrForbidden = errors.New("forbidden")
)
