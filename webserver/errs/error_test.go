package errs_test

import (
	"errors"
	"github.com/PaluMacil/dwn/webserver/errs"
	"github.com/dgraph-io/badger/v3"
	"net/http"
	"testing"
)

func TestUnwrap(t *testing.T) {
	myErr := errs.StatusError{
		Code: http.StatusConflict,
		Err:  badger.ErrConflict,
	}
	is := errors.Is(myErr, badger.ErrConflict)
	if !is {
		t.Errorf("expected Unwrapped error to be a badger.ErrConflict, but it wasn't")
	}
}
