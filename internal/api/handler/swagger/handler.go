package swaggerhandler

import (
	"context"
	"encoding/json"
	"net/http"

	api "github.com/qreaqtor/wallet/internal/generated/api"
	"github.com/qreaqtor/wallet/pkg/api/handler/request"
	"github.com/qreaqtor/wallet/pkg/api/handler/response"
)

const (
	Method = http.MethodGet

	Path = "/docs"
)

type handler struct{}

func New() *handler {
	return &handler{}
}

func (h *handler) Handle(ctx context.Context, req request.Request[any], resp response.Success[any]) error {
	swagger, err := api.GetSwagger()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		return err
	}

	resp.Raw(data)

	return nil
}
