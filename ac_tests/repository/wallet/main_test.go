//go:build integration

package wallet

import (
	"testing"

	"github.com/google/uuid"
	"github.com/qreaqtor/wallet/internal/domain/errs"
	"github.com/qreaqtor/wallet/internal/domain/types"
	"github.com/qreaqtor/wallet/internal/infrastucture/di"
	"github.com/qreaqtor/wallet/internal/infrastucture/repository/wallet"
	"github.com/qreaqtor/wallet/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	ctx := t.Context()

	cfg, err := di.LoadConfig()
	assert.NoError(t, err)

	db, err := di.NewDB(t.Context(), cfg.Database)
	assert.NoError(t, err)

	repo := wallet.New(db)

	cleanup := test.GetCleanup(ctx, db, "wallets")

	tests := map[string]func(t *testing.T){
		"success: upsert": func(t *testing.T) {
			t.Cleanup(cleanup)

			walletID := types.WalletID(uuid.New())

			wallet, err := repo.Upsert(ctx, walletID, 100)
			assert.NoError(t, err)
			assert.EqualValues(t, wallet.Balance, 100)

			wallet, err = repo.Upsert(ctx, walletID, 500)
			assert.NoError(t, err)
			assert.EqualValues(t, wallet.Balance, 500)

			_, err = repo.Upsert(ctx, walletID, -1000)
			assert.Error(t, err)
		},

		"success: GetByID": func(t *testing.T) {
			t.Cleanup(cleanup)

			walletID := types.WalletID(uuid.New())

			expected, err := repo.Upsert(ctx, walletID, 100)
			assert.NoError(t, err)

			wallet, err := repo.GetByID(ctx, walletID)
			assert.NoError(t, err)

			assert.Equal(t, expected, wallet)
		},

		"fail: upsert: check constraint": func(t *testing.T) {
			t.Cleanup(cleanup)

			walletID := types.WalletID(uuid.New())

			_, err := repo.Upsert(ctx, walletID, -1000)
			assert.Error(t, err)
		},

		"fail: GetByID: wallet not found": func(t *testing.T) {
			t.Cleanup(cleanup)

			walletID := types.WalletID(uuid.New())

			_, err := repo.GetByID(ctx, walletID)
			assert.True(t, errs.IsNotFound(err))
		},
	}

	for name, run := range tests {
		t.Run(name, run)
	}
}
