package getwallet

import (
	"context"
	"net/http"

	api "github.com/qreaqtor/wallet/internal/generated/api"
	"github.com/qreaqtor/wallet/pkg/api/handler/request"
	"github.com/qreaqtor/wallet/pkg/api/handler/response"
)

const (
	Method = http.MethodGet

	Path = "/v1/wallets/{WALLET_UUID}"
)

type handler struct{}

func New() *handler {
	return &handler{}
}

func (h *handler) Handle(ctx context.Context, req request.Request[request.Raw], resp response.Response[api.WalletBalanceResponse]) error {
	return nil
}
