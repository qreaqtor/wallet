package updatewallet

import (
	"context"
	"net/http"

	api "github.com/qreaqtor/wallet/internal/generated/api"
	"github.com/qreaqtor/wallet/pkg/api/handler/request"
	"github.com/qreaqtor/wallet/pkg/api/handler/response"
)

const (
	Method = http.MethodPost

	Path = "/v1/wallet"
)

type handler struct{}

func New() *handler {
	return &handler{}
}

func (h *handler) Handle(ctx context.Context, req request.Request[api.WalletOperationRequest], resp response.Response[api.WalletBalanceResponse]) error {
	return nil
}
