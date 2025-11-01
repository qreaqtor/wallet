package handler

import (
	"context"
	"net/http"

	"github.com/qreaqtor/wallet/pkg/api"
	errhandler "github.com/qreaqtor/wallet/pkg/api/handler/error_handler"
	"github.com/qreaqtor/wallet/pkg/api/handler/request"
	"github.com/qreaqtor/wallet/pkg/api/handler/response"
)

type handle[In, Out any] interface {
	Handle(context.Context, request.Request[In], response.Response[Out]) error
}

func New[In, Out any](log api.Log, inner handle[In, Out]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := request.New[In](r)
		resp := response.New[Out](w)

		if err := inner.Handle(ctx, req, resp); err != nil {
			errhandler.Handle(ctx, log, resp, err)
		}
	}
}
