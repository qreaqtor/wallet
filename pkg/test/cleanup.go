package test

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetCleanup(ctx context.Context, db *pgxpool.Pool, table string) func() {
	return func() {
		db.Exec(ctx, "DELETE FROM wallets;")
	}
}
