package ratelimiter

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/qreaqtor/wallet/internal/domain/types"
	"github.com/stretchr/testify/assert"
)

func TestLimiter_Allow(t *testing.T) {
	limit := 100

	l := New(int64(limit))

	walletID := types.WalletID(uuid.New())

	for range limit {
		err := l.Allow(walletID)
		assert.NoError(t, err)
	}

	assert.Error(t, l.Allow(walletID))
}

func TestLimiter_SeparateWallets(t *testing.T) {
	limit := 2

	l := New(int64(limit))

	walletID1 := types.WalletID(uuid.New())
	walletID2 := types.WalletID(uuid.New())

	assert.NoError(t, l.Allow(walletID1))
	assert.NoError(t, l.Allow(walletID2))

	assert.NoError(t, l.Allow(walletID1))
	assert.Error(t, l.Allow(walletID1))

	assert.NoError(t, l.Allow(walletID2))
}

func TestLimiter_Concurrent(t *testing.T) {
	limit := 1
	count := 100

	l := New(int64(limit))

	wg := new(sync.WaitGroup)
	for range count {
		wg.Go(func() {
			walletID := types.WalletID(uuid.New())

			assert.NoError(t, l.Allow(walletID))
			assert.Error(t, l.Allow(walletID))
		})
	}

	wg.Wait()
}
