package errs

import (
	"errors"
	"net/http"
)

// The following error handling is inspired by Matt Silverlock's blog, with the switch
// statement on the error type being used in the handler wrapper. I have added to the
// original example with additional error types.
// TODO: add validation error type
// https://blog.questionable.services/article/http-handler-error-handling-revisited/
// His snippets are MIT licensed. See bottom of this file for license text.

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

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

/*
MIT License

Copyright (c) 2019 Matt Silverlock

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

var (
	StatusUnauthorized = StatusError{http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized))}
	StatusForbidden    = StatusError{http.StatusForbidden, errors.New(http.StatusText(http.StatusForbidden))}
	StatusNotFound     = StatusError{http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound))}
	StatusLocked       = StatusError{http.StatusLocked, errors.New(http.StatusText(http.StatusLocked))}
)
