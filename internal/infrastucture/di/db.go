package di

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/qreaqtor/wallet/internal/infrastucture/di/config"
)

func NewDB(ctx context.Context, cfg config.Database) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Address,
		cfg.DatabaseName,
	)

	conn, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
