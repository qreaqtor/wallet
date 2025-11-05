package singleflight

import "context"

type Group[Out any] interface {
	Do(ctx context.Context, key string, fn func() (Out, error)) (Out, error)
}
