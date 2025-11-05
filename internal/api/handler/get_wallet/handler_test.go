package getwallet

import (
	"testing"

	"github.com/google/uuid"
	mock_getwallet "github.com/qreaqtor/wallet/internal/api/handler/get_wallet/mock"
	mock_limiter "github.com/qreaqtor/wallet/internal/api/rate_limiter/mock"
	"github.com/qreaqtor/wallet/internal/domain/entity"
	"github.com/qreaqtor/wallet/internal/domain/errs"
	"github.com/qreaqtor/wallet/internal/domain/types"
	api "github.com/qreaqtor/wallet/internal/generated/api"
	mock_request "github.com/qreaqtor/wallet/pkg/api/handler/request/mock"
	mock_response "github.com/qreaqtor/wallet/pkg/api/handler/response/mock"
	"github.com/qreaqtor/wallet/pkg/test"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandle(t *testing.T) {
	t.Parallel()
	
	ctx := t.Context()

	ctrl := gomock.NewController(t)

	type _mocks struct {
		limiter  *mock_limiter.MockLimiter
		repo     *mock_getwallet.MockwalletRepo
		request  *mock_request.MockRequest[any]
		response *mock_response.MockSuccess[api.WalletBalanceResponse]
	}

	setupRequestGetPath := func(m _mocks, id types.WalletID, err error) {
		mockParams := mock_request.NewMockParams(ctrl)

		m.request.EXPECT().GetPath().Return(mockParams)

		mockParams.EXPECT().Get(walletIDName).Return(id.UUID().String(), err)
	}

	tests := map[string]struct {
		setup   func(_mocks)
		wantErr assert.ErrorAssertionFunc
	}{
		"success": {
			setup: func(m _mocks) {
				id := types.WalletID(uuid.New())

				setupRequestGetPath(m, id, nil)

				m.limiter.EXPECT().Allow(id).Return(nil)

				m.repo.EXPECT().GetByID(ctx, id).Return(entity.Wallet{
					ID:      id,
					Balance: 100,
				}, nil)

				m.response.EXPECT().OK(api.WalletBalanceResponse{
					Balance:  100,
					WalletId: id.UUID(),
				})
			},
			wantErr: assert.NoError,
		},
		"path_get_error": {
			setup: func(m _mocks) {
				id := types.WalletID(uuid.New())

				setupRequestGetPath(m, id, assert.AnError)
			},
			wantErr: test.ErrorAs(errs.NewBadRequestErr(assert.AnError)),
		},

		"invalid_wallet_id": {
			setup: func(m _mocks) {
				mockParams := mock_request.NewMockParams(ctrl)

				m.request.EXPECT().GetPath().Return(mockParams)

				mockParams.EXPECT().Get(walletIDName).Return("invalid-uuid", nil)
			},
			wantErr: test.ErrorAs(errs.NewBadRequestErr(assert.AnError)),
		},

		"limiter_error": {
			setup: func(m _mocks) {
				id := types.WalletID(uuid.New())

				setupRequestGetPath(m, id, nil)

				m.limiter.EXPECT().Allow(id).Return(assert.AnError)
			},
			wantErr: assert.Error,
		},

		"repo_error": {
			setup: func(m _mocks) {
				id := types.WalletID(uuid.New())

				setupRequestGetPath(m, id, nil)

				m.limiter.EXPECT().Allow(id).Return(nil)

				m.repo.EXPECT().GetByID(ctx, id).Return(entity.Wallet{}, assert.AnError)
			},
			wantErr: assert.Error,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := _mocks{
				limiter:  mock_limiter.NewMockLimiter(ctrl),
				repo:     mock_getwallet.NewMockwalletRepo(ctrl),
				request:  mock_request.NewMockRequest[any](ctrl),
				response: mock_response.NewMockSuccess[api.WalletBalanceResponse](ctrl),
			}

			tc.setup(m)

			tc.wantErr(t, New(m.limiter, m.repo).Handle(ctx, m.request, m.response))
		})
	}
}
