package errs

import (
	"errors"
	"fmt"
)

type badRequestErr struct {
	cause error
}

func NewBadRequestErr(err error) *badRequestErr {
	return &badRequestErr{
		cause: fmt.Errorf("bad request: %w", err),
	}
}

func (e badRequestErr) Error() string {
	return e.cause.Error()
}

func IsBadRequest(err error) bool {
	errBadRequest := new(badRequestErr)

	return errors.As(err, &errBadRequest)
}
