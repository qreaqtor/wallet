package entity

import (
	"time"

	"github.com/qreaqtor/wallet/internal/domain/types"
)

type Wallet struct {
	ID        types.WalletID
	Balance   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
