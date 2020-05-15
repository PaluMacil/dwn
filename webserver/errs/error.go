package errs

import (
	"errors"
	"net/http"
	"strings"
)

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}

func (se StatusError) Unwrap() error {
	return se.Err
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

var (
	StatusUnauthorized = StatusError{http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized))}
	StatusForbidden    = StatusError{http.StatusForbidden, errors.New(http.StatusText(http.StatusForbidden))}
	StatusNotFound     = StatusError{http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound))}
	StatusLocked       = StatusError{http.StatusLocked, errors.New(http.StatusText(http.StatusLocked))}
)

func StatusBadRequest(problems ...string) StatusError {
	errorString := strings.Join(problems, ";")
	return StatusError{
		Code: http.StatusBadRequest,
		Err:  errors.New(errorString),
	}
}
