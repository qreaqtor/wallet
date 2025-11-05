package updatewallet

import (
	"testing"

	"github.com/google/uuid"
	mock_updatewallet "github.com/qreaqtor/wallet/internal/api/handler/update_wallet/mock"
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
		usecase  *mock_updatewallet.Mockusecase
		request  *mock_request.MockRequest[api.UpdateWalletBalanceJSONRequestBody]
		response *mock_response.MockSuccess[api.WalletBalanceResponse]
	}

	tests := map[string]struct {
		input   api.UpdateWalletBalanceJSONRequestBody
		setup   func(_mocks, api.UpdateWalletBalanceJSONRequestBody)
		wantErr assert.ErrorAssertionFunc
	}{
		"success: deposit": {
			input: api.WalletOperationRequest{
				Amount:        100,
				OperationType: api.DEPOSIT,
				WalletId:      uuid.New(),
			},
			setup: func(m _mocks, in api.UpdateWalletBalanceJSONRequestBody) {
				m.request.EXPECT().GetBody().Return(in, nil)

				id := types.WalletID(in.WalletId)
				m.limiter.EXPECT().Allow(id).Return(nil)

				action := types.Action(in.OperationType)
				amount := types.Amount(in.Amount)
				m.usecase.EXPECT().Run(ctx, id, action, amount).Return(entity.Wallet{
					ID:      id,
					Balance: in.Amount,
				}, nil)

				m.response.EXPECT().OK(api.WalletBalanceResponse{
					Balance:  in.Amount,
					WalletId: id.UUID(),
				})
			},
			wantErr: assert.NoError,
		},

		"success: withdraw": {
			input: api.WalletOperationRequest{
				Amount:        100,
				OperationType: api.WITHDRAW,
				WalletId:      uuid.New(),
			},
			setup: func(m _mocks, in api.UpdateWalletBalanceJSONRequestBody) {
				m.request.EXPECT().GetBody().Return(in, nil)

				id := types.WalletID(in.WalletId)
				m.limiter.EXPECT().Allow(id).Return(nil)

				action := types.Action(in.OperationType)
				amount := types.Amount(in.Amount)
				m.usecase.EXPECT().Run(ctx, id, action, amount).Return(entity.Wallet{
					ID:      id,
					Balance: in.Amount,
				}, nil)

				m.response.EXPECT().OK(api.WalletBalanceResponse{
					Balance:  in.Amount,
					WalletId: id.UUID(),
				})
			},
			wantErr: assert.NoError,
		},

		"body_get_error": {
			setup: func(m _mocks, in api.UpdateWalletBalanceJSONRequestBody) {
				m.request.EXPECT().GetBody().Return(api.UpdateWalletBalanceJSONRequestBody{}, assert.AnError)
			},
			wantErr: test.ErrorAs(errs.NewBadRequestErr(assert.AnError)),
		},

		"invalid_amount": {
			input: api.WalletOperationRequest{
				WalletId: uuid.New(),
				Amount:   -10,
			},
			setup: func(m _mocks, in api.UpdateWalletBalanceJSONRequestBody) {
				m.request.EXPECT().GetBody().Return(in, nil)
			},
			wantErr: test.ErrorAs(errs.NewBadRequestErr(assert.AnError)),
		},

		"invalid_wallet_id": {
			input: api.WalletOperationRequest{
				WalletId: uuid.Nil,
				Amount:   100,
			},
			setup: func(m _mocks, in api.UpdateWalletBalanceJSONRequestBody) {
				m.request.EXPECT().GetBody().Return(in, nil)
			},
			wantErr: test.ErrorAs(errs.NewBadRequestErr(assert.AnError)),
		},

		"limiter_error": {
			input: api.WalletOperationRequest{
				WalletId: uuid.New(),
				Amount:   100,
			},
			setup: func(m _mocks, in api.UpdateWalletBalanceJSONRequestBody) {
				m.request.EXPECT().GetBody().Return(in, nil)

				id := types.WalletID(in.WalletId)
				m.limiter.EXPECT().Allow(id).Return(assert.AnError)
			},
			wantErr: assert.Error,
		},

		"usecase_error": {
			input: api.WalletOperationRequest{
				Amount:        100,
				OperationType: api.DEPOSIT,
				WalletId:      uuid.New(),
			},
			setup: func(m _mocks, in api.UpdateWalletBalanceJSONRequestBody) {
				m.request.EXPECT().GetBody().Return(in, nil)

				id := types.WalletID(in.WalletId)
				m.limiter.EXPECT().Allow(id).Return(nil)

				action := types.Action(in.OperationType)
				amount := types.Amount(in.Amount)
				m.usecase.EXPECT().Run(ctx, id, action, amount).Return(entity.Wallet{}, assert.AnError)
			},
			wantErr: assert.Error,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			
			m := _mocks{
				limiter:  mock_limiter.NewMockLimiter(ctrl),
				usecase:  mock_updatewallet.NewMockusecase(ctrl),
				request:  mock_request.NewMockRequest[api.UpdateWalletBalanceJSONRequestBody](ctrl),
				response: mock_response.NewMockSuccess[api.WalletBalanceResponse](ctrl),
			}

			tc.setup(m, tc.input)

			tc.wantErr(t, New(m.limiter, m.usecase).Handle(ctx, m.request, m.response))
		})
	}
}
