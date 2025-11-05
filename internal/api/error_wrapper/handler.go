package errwrapper

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/qreaqtor/wallet/internal/domain/errs"
	"github.com/qreaqtor/wallet/pkg/api/handler/request"
	"github.com/qreaqtor/wallet/pkg/api/handler/response"
	"github.com/qreaqtor/wallet/pkg/logger"
)

func New[In, Out any](log logger.Log, inner handler[In, Out]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := request.New[In](r)
		resp := response.New[Out](w)

		if err := inner.Handle(ctx, req, resp); err != nil {
			handle(ctx, log, resp, err)
		}
	}
}

func handle(ctx context.Context, log logger.Log, resp response.Error, err error) {
	writeLog := log.Error

	switch {
	case errs.IsNotFound(err):
		resp.NotFound(err)
		writeLog = log.Info

	case isBadRequest(err):
		resp.BadRequest(err)
		writeLog = log.Debug

	case errs.IsTooManyRequests(err):
		resp.TooManyRequests(err)
		writeLog = log.Warn

	default:
		resp.InternalError(err)
	}

	writeLog(ctx, err.Error())
}

func isBadRequest(err error) bool {
	errPgConn := new(pgconn.PgError)

	return errs.IsBadRequest(err) ||
		errors.As(err, &errPgConn) && errPgConn.Code == "23514"
}
