package getwallet

import (
	"context"

	"github.com/qreaqtor/wallet/internal/domain/entity"
	"github.com/qreaqtor/wallet/internal/domain/types"
)

//go:generate mockgen -source=contract.go -destination=mock/contract.go -package=mock_getwallet

type walletRepo interface {
	GetByID(ctx context.Context, walletID types.WalletID) (entity.Wallet, error)
}
