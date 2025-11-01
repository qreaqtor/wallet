package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	getwallet "github.com/qreaqtor/wallet/internal/api/handler/get_wallet"
	updatewallet "github.com/qreaqtor/wallet/internal/api/handler/update_wallet"
	"github.com/qreaqtor/wallet/internal/config"
	"github.com/qreaqtor/wallet/pkg/api/handler"
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
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log := logger.New(cfg.Log.Level, cfg.Log.Pretty)

	apiServer := server.New(log, 8080)

	apiServer.Handle(getwallet.Method, getwallet.Path, handler.New(log, getwallet.New()))

	apiServer.Handle(updatewallet.Method, updatewallet.Path, handler.New(log, updatewallet.New()))

	return apiServer.Run(ctx)
}
