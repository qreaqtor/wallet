package wallet

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/google/uuid"
	"github.com/qreaqtor/wallet/internal/domain/entity"
	"github.com/qreaqtor/wallet/internal/domain/errs"
	"github.com/qreaqtor/wallet/internal/domain/types"
)

type (
	walletRepo struct {
		db *pgxpool.Pool
	}

	wallet struct {
		ID        uuid.UUID
		Balance   int64
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

func New(db *pgxpool.Pool) *walletRepo {
	return &walletRepo{
		db: db,
	}
}

func (r *walletRepo) GetByID(ctx context.Context, walletID types.WalletID) (entity.Wallet, error) {
	sql := `
		SELECT id, balance, created_at, updated_at FROM wallets
		WHERE id = $1;
	`

	wallet := wallet{}

	err := r.db.QueryRow(ctx, sql, walletID).Scan(
		&wallet.ID,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Wallet{}, errs.NewNotFoundErr(err)
	}
	if err != nil {
		return entity.Wallet{}, err
	}

	return toDomainWallet(wallet), nil
}

func (r *walletRepo) Upsert(ctx context.Context, walletID types.WalletID, balance int64) (entity.Wallet, error) {
	sql := `
		INSERT INTO wallets (id, balance)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET 
    		balance = EXCLUDED.balance,
    		updated_at = NOW()
		RETURNING id, balance, created_at, updated_at;
	`

	wallet := wallet{}

	err := r.db.QueryRow(ctx, sql, walletID, balance).Scan(
		&wallet.ID,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)
	if err != nil {
		return entity.Wallet{}, err
	}

	return toDomainWallet(wallet), nil
}

func toDomainWallet(wallet wallet) entity.Wallet {
	return entity.Wallet{
		ID:        types.WalletID(wallet.ID),
		Balance:   wallet.Balance,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}
}
