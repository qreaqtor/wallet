package updatewallet

import (
	"testing"

	"github.com/google/uuid"
	"github.com/qreaqtor/wallet/internal/domain/entity"
	"github.com/qreaqtor/wallet/internal/domain/errs"
	"github.com/qreaqtor/wallet/internal/domain/types"
	mock_updatewallet "github.com/qreaqtor/wallet/internal/usecase/update_wallet/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUsecase(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	ctrl := gomock.NewController(t)

	type _mocks struct {
		repo *mock_updatewallet.Mockrepo
	}

	walletID := types.WalletID(uuid.New())
	balance := int64(1000)
	amount := int64(100)

	tests := map[string]struct {
		setup    func(_mocks)
		action   types.Action
		expected entity.Wallet
		wantErr  assert.ErrorAssertionFunc
	}{
		"success: deposit: wallet exists": {
			action: types.ActionDeposit,
			setup: func(m _mocks) {
				m.repo.EXPECT().GetByID(ctx, walletID).Return(entity.Wallet{
					ID:      walletID,
					Balance: balance,
				}, nil)

				m.repo.EXPECT().Upsert(ctx, walletID, balance+amount).Return(entity.Wallet{
					ID:      walletID,
					Balance: balance + amount,
				}, nil)
			},
			expected: entity.Wallet{
				ID:      walletID,
				Balance: balance + amount,
			},
			wantErr: assert.NoError,
		},

		"success: deposit: wallet not exists": {
			action: types.ActionDeposit,
			setup: func(m _mocks) {
				m.repo.EXPECT().GetByID(ctx, walletID).Return(entity.Wallet{}, errs.NewNotFoundErr(assert.AnError))

				m.repo.EXPECT().Upsert(ctx, walletID, amount).Return(entity.Wallet{
					ID:      walletID,
					Balance: amount,
				}, nil)
			},
			expected: entity.Wallet{
				ID:      walletID,
				Balance: amount,
			},
			wantErr: assert.NoError,
		},

		"deposit: repo getByID error": {
			action: types.ActionDeposit,
			setup: func(m _mocks) {
				m.repo.EXPECT().GetByID(ctx, walletID).Return(entity.Wallet{}, assert.AnError)
			},
			wantErr: assert.Error,
		},

		"success: withdraw": {
			action: types.ActionWITHDRAW,
			setup: func(m _mocks) {
				m.repo.EXPECT().GetByID(ctx, walletID).Return(entity.Wallet{
					ID:      walletID,
					Balance: balance,
				}, nil)

				m.repo.EXPECT().Upsert(ctx, walletID, balance-amount).Return(entity.Wallet{
					ID:      walletID,
					Balance: balance - amount,
				}, nil)
			},
			expected: entity.Wallet{
				ID:      walletID,
				Balance: balance - amount,
			},
			wantErr: assert.NoError,
		},

		"withdraw: repo getByID error": {
			action: types.ActionWITHDRAW,
			setup: func(m _mocks) {
				m.repo.EXPECT().GetByID(ctx, walletID).Return(entity.Wallet{}, assert.AnError)
			},
			wantErr: assert.Error,
		},

		"withdraw: insufficient balance error": {
			action: types.ActionWITHDRAW,
			setup: func(m _mocks) {
				m.repo.EXPECT().GetByID(ctx, walletID).Return(entity.Wallet{
					ID:      walletID,
					Balance: 1,
				}, nil)
			},
			wantErr: assert.Error,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			
			m := _mocks{
				repo: mock_updatewallet.NewMockrepo(ctrl),
			}

			tc.setup(m)

			wallet, err := New(m.repo).Run(ctx, walletID, tc.action, types.Amount(amount))
			tc.wantErr(t, err)

			assert.Equal(t, tc.expected, wallet)
		})
	}
}
