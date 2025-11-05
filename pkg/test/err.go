package test

import (
	"github.com/stretchr/testify/assert"
)

func ErrorAs[E error](target E) func(tt assert.TestingT, err error, args ...any) bool {
	return func(tt assert.TestingT, err error, args ...any) bool {
		return assert.ErrorAs(tt, err, &target, args)
	}
}
