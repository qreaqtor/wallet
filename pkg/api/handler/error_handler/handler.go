package errhandler

import (
	"context"
	"errors"

	"github.com/qreaqtor/wallet/pkg/api"
	"github.com/qreaqtor/wallet/pkg/api/handler/response"
)

func Handle(ctx context.Context, log api.Log, resp response.ErrorResponse, err error) {
	if err == nil {
		return
	}

	writeLog := log.Error

	switch {
	case errors.As(err, new(notFoundErr)):
		resp.NotFound(err)
		writeLog = log.Info

	case errors.As(err, new(badRequestErr)):
		resp.BadRequest(err)
		writeLog = log.Warn

	default:
		resp.InternalError(err)
	}

	writeLog(ctx, err.Error())
}
