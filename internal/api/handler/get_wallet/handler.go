package getwallet

import (
	"context"
	"net/http"

	ratelimiter "github.com/qreaqtor/wallet/internal/api/rate_limiter"
	"github.com/qreaqtor/wallet/internal/domain/entity"
	"github.com/qreaqtor/wallet/internal/domain/errs"
	"github.com/qreaqtor/wallet/internal/domain/types"
	api "github.com/qreaqtor/wallet/internal/generated/api"
	"github.com/qreaqtor/wallet/pkg/api/handler/request"
	"github.com/qreaqtor/wallet/pkg/api/handler/response"
)

const (
	Method = http.MethodGet

	walletIDName = "WALLET_UUID"

	Path = "/v1/wallets/{" + walletIDName + "}"
)

type handler struct {
	repo    walletRepo
	limiter ratelimiter.Limiter
}

func New(limiter ratelimiter.Limiter, repo walletRepo) *handler {
	return &handler{
		repo:    repo,
		limiter: limiter,
	}
}

func (h *handler) Handle(ctx context.Context, req request.Request[any], resp response.Success[api.WalletBalanceResponse]) error {
	path := req.GetPath()

	rawWalletID, err := path.Get(walletIDName)
	if err != nil {
		return errs.NewBadRequestErr(err)
	}

	walletID, err := types.NewWalletID(rawWalletID)
	if err != nil {
		return err
	}

	if err := h.limiter.Allow(walletID); err != nil {
		return err
	}

	out, err := h.repo.GetByID(ctx, walletID)
	if err != nil {
		return err
	}

	resp.OK(toWalletBalanceResponse(out))

	return nil
}

func toWalletBalanceResponse(wallet entity.Wallet) api.WalletBalanceResponse {
	return api.WalletBalanceResponse{
		WalletId:  wallet.ID.UUID(),
		Balance:   wallet.Balance,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}
}
