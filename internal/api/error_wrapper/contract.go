package errwrapper

import (
	"context"

	"github.com/qreaqtor/wallet/pkg/api/handler/request"
	"github.com/qreaqtor/wallet/pkg/api/handler/response"
)

type handler[In, Out any] interface {
	Handle(context.Context, request.Request[In], response.Success[Out]) error
}
