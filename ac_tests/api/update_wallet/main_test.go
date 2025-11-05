//go:build integration

package updatewallet

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

func Test_UpdateWallet_Success(t *testing.T) {
	ctx := t.Context()

	cfg, err := di.LoadConfig()
	assert.NoError(t, err)

	db, err := di.NewDB(ctx, cfg.Database)
	assert.NoError(t, err)

	cleanup := test.GetCleanup(ctx, db, "wallets")

	repo := wallet.New(db)

	addr := fmt.Sprintf("http://localhost:%d", cfg.Port)

	client, err := api.NewClient(addr)
	assert.NoError(t, err)

	tests := map[string]struct {
		status  int
		req     api.UpdateWalletBalanceJSONRequestBody
		getResp func(uuid.UUID) api.WalletBalanceResponse
		setup   func(types.WalletID)
	}{
		"deposit": {
			status: http.StatusOK,
			req: api.WalletOperationRequest{
				Amount:        100,
				OperationType: api.DEPOSIT,
				WalletId:      uuid.New(),
			},
			getResp: func(id uuid.UUID) api.WalletBalanceResponse {
				return api.WalletBalanceResponse{
					Balance:  100,
					WalletId: id,
				}
			},
		},

		"withdraw": {
			status: http.StatusOK,
			req: api.WalletOperationRequest{
				Amount:        100,
				OperationType: api.WITHDRAW,
				WalletId:      uuid.New(),
			},
			getResp: func(id uuid.UUID) api.WalletBalanceResponse {
				return api.WalletBalanceResponse{
					Balance:  900,
					WalletId: id,
				}
			},
			setup: func(id types.WalletID) {
				_, err := repo.Upsert(ctx, id, 1000)
				assert.NoError(t, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup)

			if tc.setup != nil {
				tc.setup(types.WalletID(tc.req.WalletId))
			}

			resp, err := client.UpdateWalletBalance(ctx, tc.req)
			assert.NoError(t, err)

			assert.Equal(t, tc.status, resp.StatusCode)

			actual := api.WalletBalanceResponse{}
			expected := tc.getResp(tc.req.WalletId)

			assert.NoError(t, test.ParseRespBody(resp, &actual))

			actual.CreatedAt = expected.CreatedAt
			actual.UpdatedAt = expected.UpdatedAt

			assert.Equal(t, expected, actual)
		})
	}
}

func Test_UpdateWallet_Error(t *testing.T) {
	ctx := t.Context()

	cfg, err := di.LoadConfig()
	assert.NoError(t, err)

	db, err := di.NewDB(ctx, cfg.Database)
	assert.NoError(t, err)

	cleanup := test.GetCleanup(ctx, db, "wallets")

	repo := wallet.New(db)

	addr := fmt.Sprintf("http://localhost:%d", cfg.Port)

	client, err := api.NewClient(addr)
	assert.NoError(t, err)

	tests := map[string]struct {
		status int
		req    api.UpdateWalletBalanceJSONRequestBody
		setup  func(types.WalletID)
	}{
		"deposit: negative amount": {
			status: http.StatusBadRequest,
			req: api.WalletOperationRequest{
				Amount:        -100,
				OperationType: api.DEPOSIT,
				WalletId:      uuid.New(),
			},
		},

		"withdraw: balance less than amount": {
			status: http.StatusBadRequest,
			req: api.WalletOperationRequest{
				Amount:        200,
				OperationType: api.WITHDRAW,
				WalletId:      uuid.New(),
			},
			setup: func(id types.WalletID) {
				_, err := repo.Upsert(ctx, id, 100)
				assert.NoError(t, err)
			},
		},

		"withdraw: not found": {
			status: http.StatusNotFound,
			req: api.WalletOperationRequest{
				Amount:        200,
				OperationType: api.WITHDRAW,
				WalletId:      uuid.New(),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup)

			if tc.setup != nil {
				tc.setup(types.WalletID(tc.req.WalletId))
			}

			resp, err := client.UpdateWalletBalance(ctx, tc.req)
			assert.NoError(t, err)

			assert.Equal(t, tc.status, resp.StatusCode)
		})
	}
}
