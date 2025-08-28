package cache

import (
	"context"
	"fmt"
	"github.com/Killazius/L0/internal/domain"
	"golang.org/x/sync/errgroup"
)

type OrderProvider interface {
	GetAll(ctx context.Context) ([]domain.Order, error)
}

func (c *Cache) Restore(ctx context.Context, repo OrderProvider, workers int) error {
	orders, err := repo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all orders: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(workers)

	for _, order := range orders {
		order := order
		g.Go(func() error {
			return c.Set(ctx, &order)
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to restore orders: %w", err)
	}

	return nil
}
