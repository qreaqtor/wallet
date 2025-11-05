package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	errwrapper "github.com/qreaqtor/wallet/internal/api/error_wrapper"
	getwallet_h "github.com/qreaqtor/wallet/internal/api/handler/get_wallet"
	swaggerhandler "github.com/qreaqtor/wallet/internal/api/handler/swagger"
	updatewallet_h "github.com/qreaqtor/wallet/internal/api/handler/update_wallet"
	ratelimiter "github.com/qreaqtor/wallet/internal/api/rate_limiter"
	wallet_cache "github.com/qreaqtor/wallet/internal/infrastucture/cache/wallet"
	"github.com/qreaqtor/wallet/internal/infrastucture/di"
	wallet_repo "github.com/qreaqtor/wallet/internal/infrastucture/repository/wallet"
	updatewallet_uc "github.com/qreaqtor/wallet/internal/usecase/update_wallet"
	"github.com/qreaqtor/wallet/pkg/api/server"
	"github.com/qreaqtor/wallet/pkg/logger"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx); err != nil {
		log.Fatalln(err)
	}
}

func run(ctx context.Context) error {
	cfg, err := di.LoadConfig()
	if err != nil {
		return err
	}

	db, err := di.NewDB(ctx, cfg.Database)
	if err != nil {
		return err
	}

	log := logger.New(cfg.Log.Level, cfg.Log.Pretty)

	limiter := ratelimiter.New(cfg.LimiterRPS)

	walletRepo := wallet_repo.New(db)

	walletCache := wallet_cache.New(walletRepo)

	updateUC := updatewallet_uc.New(walletCache)

	apiServer := server.New(log, cfg.Port)

	apiServer.Handle(
		swaggerhandler.Method,
		swaggerhandler.Path,
		errwrapper.New(log, swaggerhandler.New()),
	)

	apiServer.Handle(
		getwallet_h.Method,
		getwallet_h.Path,
		errwrapper.New(log, getwallet_h.New(limiter, walletCache)),
	)

	apiServer.Handle(
		updatewallet_h.Method,
		updatewallet_h.Path,
		errwrapper.New(log, updatewallet_h.New(limiter, updateUC)),
	)

	return apiServer.Run(ctx)
}
