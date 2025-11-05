package updatewallet

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
	Method = http.MethodPost

	Path = "/v1/wallet"
)

type handler struct {
	uc usecase

	limiter ratelimiter.Limiter
}

func New(limiter ratelimiter.Limiter, uc usecase) *handler {
	return &handler{
		uc:      uc,
		limiter: limiter,
	}
}

func (h *handler) Handle(ctx context.Context, req request.Request[api.UpdateWalletBalanceJSONRequestBody], resp response.Success[api.WalletBalanceResponse]) error {
	body, err := req.GetBody()
	if err != nil {
		return errs.NewBadRequestErr(err)
	}

	amount, err := types.NewAmount(body.Amount)
	if err != nil {
		return err
	}

	walletID, err := types.NewWalletID(body.WalletId.String())
	if err != nil {
		return err
	}

	if err := h.limiter.Allow(walletID); err != nil {
		return err
	}

	wallet, err := h.uc.Run(ctx, walletID, types.Action(body.OperationType), amount)
	if err != nil {
		return err
	}

	resp.OK(toWalletBalanceResponse(wallet))

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
