package types

import (
	"fmt"

	"github.com/qreaqtor/wallet/internal/domain/errs"
)

type Amount uint32

func NewAmount(v int64) (Amount, error) {
	if v <= 0 {
		return 0, errs.NewBadRequestErr(fmt.Errorf("negative amount: %v", v))
	}

	return Amount(v), nil
}

func (a Amount) ToInt64() int64 {
	return int64(a)
}
