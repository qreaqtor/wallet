package updatewallet

import (
	"context"

	"github.com/qreaqtor/wallet/internal/domain/entity"
	"github.com/qreaqtor/wallet/internal/domain/types"
)

//go:generate mockgen -source=contract.go -destination=mock/contract.go -package=mock_updatewallet

type repo interface {
	GetByID(ctx context.Context, walletID types.WalletID) (entity.Wallet, error)
	Upsert(ctx context.Context, walletID types.WalletID, balance int64) (entity.Wallet, error)
}
