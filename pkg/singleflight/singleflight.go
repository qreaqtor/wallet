package singleflight

import (
	"context"
	"fmt"

	"golang.org/x/sync/singleflight"
)

type group[Out any] struct {
	group singleflight.Group
}

func New[Out any]() Group[Out] {
	return &group[Out]{
		group: singleflight.Group{},
	}
}

func (g *group[Out]) Do(ctx context.Context, key string, fn func() (Out, error)) (Out, error) {
	// The returned channel will not be closed.
	// See https://cs.opensource.google/go/x/sync/+/refs/tags/v0.17.0:singleflight/singleflight.go;l=120
	resultChan := g.group.DoChan(key, func() (any, error) {
		return fn()
	})

	var d Out

	select {
	case <-ctx.Done():
		return d, context.Canceled
	case v := <-resultChan:
		if v.Err != nil {
			return d, v.Err
		}

		out, ok := v.Val.(Out)
		if !ok {
			return d, fmt.Errorf("unexpected result type: %T", v.Val)
		}

		return out, nil
	}
}
