package wallet

import (
	"context"
	"sync"

	"github.com/qreaqtor/wallet/internal/domain/entity"
	"github.com/qreaqtor/wallet/internal/domain/types"
	"github.com/qreaqtor/wallet/pkg/singleflight"
)

type cache struct {
	repo repo

	data map[types.WalletID]entity.Wallet

	mu sync.RWMutex

	group singleflight.Group[entity.Wallet]
}

func New(repo repo) *cache {
	return &cache{
		data:  map[types.WalletID]entity.Wallet{},
		mu:    sync.RWMutex{},
		repo:  repo,
		group: singleflight.New[entity.Wallet](),
	}
}

func (c *cache) GetByID(ctx context.Context, walletID types.WalletID) (entity.Wallet, error) {
	c.mu.RLock()

	if v, ok := c.data[walletID]; ok {
		c.mu.RUnlock()

		return v, nil
	}

	c.mu.RUnlock()

	wallet, err := c.group.Do(ctx, walletID.String(), func() (entity.Wallet, error) {
		wallet, err := c.repo.GetByID(ctx, walletID)
		if err != nil {
			return entity.Wallet{}, err
		}

		c.mu.Lock()
		defer c.mu.Unlock()

		if v, ok := c.data[walletID]; ok {
			return v, nil
		}

		c.data[wallet.ID] = wallet

		return wallet, nil
	})
	if err != nil {
		return entity.Wallet{}, err
	}

	return wallet, nil
}

func (c *cache) Upsert(ctx context.Context, walletID types.WalletID, balance int64) (entity.Wallet, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	wallet, err := c.repo.Upsert(ctx, walletID, balance)
	if err != nil {
		return entity.Wallet{}, err
	}

	c.data[wallet.ID] = wallet

	return wallet, nil
}
