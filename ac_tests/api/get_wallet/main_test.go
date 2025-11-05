//go:build integration

package getwallet

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/qreaqtor/wallet/internal/domain/types"
	api "github.com/qreaqtor/wallet/internal/generated/api"
	"github.com/qreaqtor/wallet/internal/infrastucture/di"
	"github.com/qreaqtor/wallet/internal/infrastucture/repository/wallet"
	"github.com/qreaqtor/wallet/pkg/test"
	"github.com/stretchr/testify/assert"
)

func Test_GetWallet_Success(t *testing.T) {
	ctx := t.Context()

	cfg, err := di.LoadConfig()
	assert.NoError(t, err)

	db, err := di.NewDB(ctx, cfg.Database)
	assert.NoError(t, err)

	t.Cleanup(test.GetCleanup(ctx, db, "wallets"))

	repo := wallet.New(db)

	walletID := types.WalletID(uuid.New())
	balance := int64(1000)

	_, err = repo.Upsert(ctx, walletID, balance)
	assert.NoError(t, err)

	addr := fmt.Sprintf("http://localhost:%d", cfg.Port)

	client, err := api.NewClient(addr)
	assert.NoError(t, err)

	resp, err := client.GetWalletBalance(ctx, walletID.UUID())
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	actual := api.WalletBalanceResponse{}
	expected := api.WalletBalanceResponse{
		Balance:  balance,
		WalletId: walletID.UUID(),
	}

	assert.NoError(t, test.ParseRespBody(resp, &actual))

	actual.CreatedAt = expected.CreatedAt
	actual.UpdatedAt = expected.UpdatedAt

	assert.Equal(t, expected, actual)
}

func Test_GetWallet_Error(t *testing.T) {
	ctx := t.Context()

	cfg, err := di.LoadConfig()
	assert.NoError(t, err)

	db, err := di.NewDB(ctx, cfg.Database)
	assert.NoError(t, err)

	cleanup := test.GetCleanup(ctx, db, "wallets")

	addr := fmt.Sprintf("http://localhost:%d", cfg.Port)

	client, err := api.NewClient(addr)
	assert.NoError(t, err)

	tests := map[string]struct {
		status int
		id     uuid.UUID
	}{
		"not found": {
			id:     uuid.New(),
			status: http.StatusNotFound,
		},
		"bad request": {
			id:     uuid.Nil,
			status: http.StatusBadRequest,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup)

			resp, err := client.GetWalletBalance(ctx, tc.id)
			assert.NoError(t, err)

			assert.Equal(t, tc.status, resp.StatusCode)
		})
	}
}
