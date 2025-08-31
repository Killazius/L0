package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Killazius/L0/internal/domain"
	"github.com/Killazius/L0/internal/repository"
	"github.com/Killazius/L0/pkg/validate"
	"go.uber.org/zap"
)

func (s *Service) GetOrder(ctx context.Context, uid string) (*domain.Order, error) {
	cacheCtx, cancel := context.WithTimeout(ctx, cacheTimeout)
	defer cancel()
	order, err := s.cache.Get(cacheCtx, uid)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("cache operation canceled: %w", err)
		}
		if !errors.Is(err, repository.ErrOrderNotFound) {
			zap.L().Warn("cache error, falling back to database", zap.String("order_uid", uid), zap.Error(err))
		}

		order, err = s.repo.Get(ctx, uid)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrOrderNotFound):
				return nil, fmt.Errorf("%w: %s", ErrOrderNotFound, uid)
			case errors.Is(err, repository.ErrDeliveryNotFound),
				errors.Is(err, repository.ErrPaymentNotFound),
				errors.Is(err, repository.ErrItemsNotFound):
				return nil, fmt.Errorf("%w: incomplete data for order %s", ErrInvalidOrderData, uid)
			default:
				return nil, fmt.Errorf("failed to get order from database: %w", err)
			}
		}
		s.wg.Go(func() {
			cacheCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), cacheTimeout)
			defer cancel()
			if err := s.cache.Set(cacheCtx, order); err != nil && !errors.Is(err, context.Canceled) {
				zap.L().Warn("failed to cache order", zap.String("order_uid", uid), zap.Error(err))
			}
		})
		zap.L().Info("from database", zap.String("uid", uid))
		return order, nil
	}
	zap.L().Info("from cache", zap.String("uid", uid))
	return order, nil
}

func (s *Service) CreateOrder(ctx context.Context, order *domain.Order) error {
	if order == nil {
		return fmt.Errorf("%w: order is nil", ErrInvalidOrderData)
	}
	if err := validate.Order(order); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidOrderData, err)
	}
	err := s.repo.Create(ctx, order)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateOrder):
			return ErrOrderAlreadyExists
		default:
			return fmt.Errorf("failed to create order: %w", err)
		}
	}
	s.wg.Go(func() {
		cacheCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), cacheTimeout)
		defer cancel()
		if err := s.cache.Set(cacheCtx, order); err != nil && !errors.Is(err, context.Canceled) {
			zap.L().Warn("failed to cache order", zap.String("order_uid", order.OrderUID), zap.Error(err))
		}
	})

	return nil
}
