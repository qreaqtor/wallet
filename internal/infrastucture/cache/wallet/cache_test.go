package wallet

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/qreaqtor/wallet/internal/domain/entity"
	"github.com/qreaqtor/wallet/internal/domain/types"
	mock_cache "github.com/qreaqtor/wallet/internal/infrastucture/cache/wallet/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_Upsert_Success(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	ctrl := gomock.NewController(t)

	mockRepo := mock_cache.NewMockrepo(ctrl)

	cache := New(mockRepo)

	expected := entity.Wallet{
		ID:      types.WalletID(uuid.New()),
		Balance: int64(100),
	}

	mockRepo.EXPECT().Upsert(ctx, expected.ID, expected.Balance).Return(expected, nil)

	_, err := cache.Upsert(ctx, expected.ID, expected.Balance)
	assert.NoError(t, err)

	wallet, ok := cache.data[expected.ID]
	assert.True(t, ok)
	assert.Equal(t, expected, wallet)
}

func Test_Upsert_Error(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	ctrl := gomock.NewController(t)

	mockRepo := mock_cache.NewMockrepo(ctrl)

	cache := New(mockRepo)

	id := types.WalletID(uuid.New())
	balance := int64(100)

	mockRepo.EXPECT().Upsert(ctx, id, balance).Return(entity.Wallet{}, assert.AnError)

	_, err := cache.Upsert(ctx, id, balance)
	assert.Error(t, err)
}

func Test_GetByID_Cache_Success(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	ctrl := gomock.NewController(t)

	mockRepo := mock_cache.NewMockrepo(ctrl)

	cache := New(mockRepo)

	expected := entity.Wallet{
		ID:      types.WalletID(uuid.New()),
		Balance: 100,
	}

	cache.data[expected.ID] = expected

	mockRepo.EXPECT().GetByID(ctx, expected.ID).Times(0)

	wallet, err := cache.GetByID(ctx, expected.ID)
	assert.NoError(t, err)
	assert.Equal(t, expected, wallet)
}

func Test_GetByID_Repo_Success(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	ctrl := gomock.NewController(t)
	mockRepo := mock_cache.NewMockrepo(ctrl)

	cache := New(mockRepo)

	expected := entity.Wallet{
		ID:      types.WalletID(uuid.New()),
		Balance: 200,
	}

	mockRepo.EXPECT().GetByID(ctx, expected.ID).Return(expected, nil)

	wallet, err := cache.GetByID(ctx, expected.ID)
	assert.NoError(t, err)
	assert.Equal(t, expected, wallet)

	got, ok := cache.data[expected.ID]
	assert.True(t, ok)
	assert.Equal(t, expected, got)
}

func Test_GetByID_Repo_Error(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	ctrl := gomock.NewController(t)
	mockRepo := mock_cache.NewMockrepo(ctrl)

	cache := New(mockRepo)

	id := types.WalletID(uuid.New())

	mockRepo.EXPECT().GetByID(ctx, id).Return(entity.Wallet{}, assert.AnError)

	_, err := cache.GetByID(ctx, id)
	assert.Error(t, err)

	_, ok := cache.data[id]
	assert.False(t, ok)
}

func Test_Upsert_And_GetByID(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	ctrl := gomock.NewController(t)

	mockRepo := mock_cache.NewMockrepo(ctrl)

	cache := New(mockRepo)

	id := types.WalletID(uuid.New())
	beforeUpsert := entity.Wallet{ID: id, Balance: 200}
	afterUpsert := entity.Wallet{ID: id, Balance: 2000}

	canUpsert := make(chan struct{})
	canGetByID := make(chan struct{})

	mockRepo.EXPECT().
		GetByID(ctx, id).
		DoAndReturn(func(context.Context, types.WalletID) (entity.Wallet, error) {
			close(canUpsert)

			<-canGetByID

			return beforeUpsert, nil
		})

	mockRepo.EXPECT().Upsert(ctx, id, afterUpsert.Balance).Return(afterUpsert, nil)

	t.Run("getByID", func(t *testing.T) {
		t.Parallel()

		wallet, err := cache.GetByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, afterUpsert, wallet)

		got, ok := cache.data[afterUpsert.ID]
		assert.True(t, ok)
		assert.Equal(t, afterUpsert, got)
	})

	t.Run("upsert", func(t *testing.T) {
		t.Parallel()

		<-canUpsert

		wallet, err := cache.Upsert(ctx, id, afterUpsert.Balance)
		assert.NoError(t, err)
		assert.Equal(t, afterUpsert, wallet)

		got, ok := cache.data[afterUpsert.ID]
		assert.True(t, ok)
		assert.Equal(t, afterUpsert, got)

		close(canGetByID)
	})
}
