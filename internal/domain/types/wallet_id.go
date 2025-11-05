package types

import (
	"errors"

	"github.com/google/uuid"
	"github.com/qreaqtor/wallet/internal/domain/errs"
)

var errNilWalletID = errors.New("nil wallet id")

type WalletID uuid.UUID

func NewWalletID(s string) (WalletID, error) {
	walletID, err := uuid.Parse(s)
	if err != nil {
		return WalletID{}, errs.NewBadRequestErr(err)
	} else if walletID == uuid.Nil {
		return WalletID{}, errs.NewBadRequestErr(errNilWalletID)
	}

	return WalletID(walletID), nil
}

func (w WalletID) UUID() uuid.UUID {
	return uuid.UUID(w)
}

func (w WalletID) String() string {
	return uuid.UUID(w).String()
}
