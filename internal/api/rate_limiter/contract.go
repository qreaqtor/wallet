package ratelimiter

import "github.com/qreaqtor/wallet/internal/domain/types"

//go:generate mockgen -source=contract.go -destination=mock/contract.go -package=mock_limiter

type Limiter interface {
	Allow(types.WalletID) error
}
