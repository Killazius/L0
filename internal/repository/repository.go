package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Killazius/L0/internal/domain"
	"golang.org/x/sync/errgroup"
)

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrDeliveryNotFound = errors.New("delivery not found")
	ErrPaymentNotFound  = errors.New("payment not found")
	ErrItemsNotFound    = errors.New("items not found")
	ErrDuplicateOrder   = errors.New("duplicate order")
)

type OrderProvider interface {
	GetAll(ctx context.Context) ([]domain.Order, error)
}
type OrderSetter interface {
	Set(ctx context.Context, order *domain.Order) error
}

func Restore(ctx context.Context, repo OrderProvider, cache OrderSetter, workers int) error {
	orders, err := repo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all orders: %w", err)
	}
	if len(orders) == 0 {
		return nil
	}

	g, ctx := errgroup.WithContext(ctx)
	if workers > 0 {
		g.SetLimit(workers)
	}

	for _, order := range orders {
		order := order
		g.Go(func() error {
			return cache.Set(ctx, &order)
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to restore orders: %w", err)
	}

	return nil
}
