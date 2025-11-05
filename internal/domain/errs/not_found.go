package errs

import (
	"errors"
	"fmt"
)

type notFoundErr struct {
	cause error
}

func NewNotFoundErr(err error) *notFoundErr {
	return &notFoundErr{
		cause: fmt.Errorf("not found: %w", err),
	}
}

func (e notFoundErr) Error() string {
	return e.cause.Error()
}

func IsNotFound(err error) bool {
	errNotFound := new(notFoundErr)

	return errors.As(err, &errNotFound)
}
