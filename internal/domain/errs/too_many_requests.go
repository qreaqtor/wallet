package errs

import (
	"errors"
	"fmt"
)

type tooManyRequestsErr struct {
	cause string
}

func NewTooManyRequestsErr(err error) *tooManyRequestsErr {
	return &tooManyRequestsErr{
		cause: fmt.Sprintf("too many requests: %v", err),
	}
}

func (e tooManyRequestsErr) Error() string {
	return e.cause
}

func IsTooManyRequests(err error) bool {
	errTooManyRequests := new(tooManyRequestsErr)

	return errors.As(err, &errTooManyRequests)
}
