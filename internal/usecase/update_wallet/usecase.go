package updatewallet

import (
	"context"
	"errors"

	"github.com/qreaqtor/wallet/internal/domain/entity"
	"github.com/qreaqtor/wallet/internal/domain/errs"
	"github.com/qreaqtor/wallet/internal/domain/types"
)

var errInsufficientBalance = errs.NewBadRequestErr(errors.New("insufficient balance"))

type usecase struct {
	wallet repo
}

func New(wallet repo) *usecase {
	return &usecase{
		wallet: wallet,
	}
}

func (uc *usecase) Run(ctx context.Context, walletID types.WalletID, action types.Action, amount types.Amount) (entity.Wallet, error) {
	if action == types.ActionDeposit {
		return uc.proccessDeposit(ctx, walletID, amount)
	}

	return uc.proccessWithdraw(ctx, walletID, amount)
}

func (uc *usecase) proccessDeposit(ctx context.Context, walletID types.WalletID, amount types.Amount) (entity.Wallet, error) {
	wallet, err := uc.wallet.GetByID(ctx, walletID)
	if err != nil && !errs.IsNotFound(err) {
		return entity.Wallet{}, err
	}

	newBalance := wallet.Balance + amount.ToInt64()

	return uc.wallet.Upsert(ctx, walletID, newBalance)
}

func (uc *usecase) proccessWithdraw(ctx context.Context, walletID types.WalletID, amount types.Amount) (entity.Wallet, error) {
	wallet, err := uc.wallet.GetByID(ctx, walletID)
	if err != nil {
		return entity.Wallet{}, err
	}

	newBalance := wallet.Balance - amount.ToInt64()

	if newBalance < 0 {
		return entity.Wallet{}, errInsufficientBalance
	}

	return uc.wallet.Upsert(ctx, walletID, newBalance)
}
