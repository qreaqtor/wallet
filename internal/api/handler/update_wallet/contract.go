package updatewallet

import (
	"context"

	"github.com/qreaqtor/wallet/internal/domain/entity"
	"github.com/qreaqtor/wallet/internal/domain/types"
)

//go:generate mockgen -source=contract.go -destination=mock/contract.go -package=mock_updatewallet

type usecase interface {
	Run(ctx context.Context, walletID types.WalletID, action types.Action, amount types.Amount) (entity.Wallet, error)
}
