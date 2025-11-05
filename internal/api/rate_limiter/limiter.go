package ratelimiter

import (
	"errors"
	"sync"

	"github.com/qreaqtor/wallet/internal/domain/errs"
	"github.com/qreaqtor/wallet/internal/domain/types"
	"golang.org/x/time/rate"
)

var errLimitReached = errors.New("limit has been reached, try later")

type limiter struct {
	items map[types.WalletID]*rate.Limiter

	mu sync.RWMutex

	rps   rate.Limit
	burst int
}

func New(rps int64) *limiter {
	return &limiter{
		items: map[types.WalletID]*rate.Limiter{},
		mu:    sync.RWMutex{},
		rps:   rate.Limit(rps),
		burst: int(rps),
	}
}

func (l *limiter) Allow(walletID types.WalletID) error {
	lim := l.getLimiter(walletID)

	if lim.Allow() {
		return nil
	}

	return errs.NewTooManyRequestsErr(errLimitReached)
}

func (l *limiter) getLimiter(walletID types.WalletID) *rate.Limiter {
	l.mu.RLock()
	if v, ok := l.items[walletID]; ok {
		l.mu.RUnlock()

		return v
	}
	l.mu.RUnlock()

	l.mu.Lock()
	defer l.mu.Unlock()

	lim := rate.NewLimiter(l.rps, l.burst)

	l.items[walletID] = lim

	return lim
}
