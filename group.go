package fanin

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Group[T any] struct {
	eg      *errgroup.Group
	mu      sync.Mutex
	results []T
}

func WithContext[T any](ctx context.Context, capacity int) (*Group[T], context.Context) {
	eg, ctx := errgroup.WithContext(ctx)
	return &Group[T]{
		eg:      eg,
		results: make([]T, 0, capacity),
	}, ctx
}

func (g *Group[T]) Go(fn func() (*T, error)) {
	g.eg.Go(func() error {
		t, err := fn()
		if err != nil {
			return err
		}
		g.mu.Lock()
		defer g.mu.Unlock()
		g.results = append(g.results, *t)
		return nil
	})
}

func (g *Group[T]) Wait() ([]T, error) {
	err := g.eg.Wait()
	if err != nil {
		return nil, err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.results, nil
}
