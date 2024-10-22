package fanin

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Group[T any] struct {
	eg      *errgroup.Group
	results chan T
}

func WithContext[T any](ctx context.Context) (*Group[T], context.Context) {
	eg, ctx := errgroup.WithContext(ctx)
	return &Group[T]{
		eg:      eg,
		results: make(chan T),
	}, ctx
}

func (g *Group[T]) Go(fn func() (*T, error)) {
	g.eg.Go(func() error {
		t, err := fn()
		if err != nil {
			return err
		}
		g.results <- *t
		return nil
	})
}

func (g *Group[T]) Wait() ([]T, error) {
	var results []T
	go func() {
		for result := range g.results {
			results = append(results, result)
		}
	}()
	err := g.eg.Wait()
	close(g.results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
