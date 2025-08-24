package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Killazius/L0/internal/domain"
	"github.com/Killazius/L0/internal/repository"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrInvalidOrderData   = errors.New("invalid order data")
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}
func (s *Service) GetOrder(ctx context.Context, uid string) (*domain.Order, error) {
	order, err := s.repo.Get(ctx, uid)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrOrderNotFound):
			return nil, fmt.Errorf("%w: %s", ErrOrderNotFound, uid)
		case errors.Is(err, repository.ErrDeliveryNotFound),
			errors.Is(err, repository.ErrPaymentNotFound),
			errors.Is(err, repository.ErrItemsNotFound):
			return nil, fmt.Errorf("%w: incomplete data for order %s", ErrInvalidOrderData, uid)
		default:
			return nil, fmt.Errorf("failed to get order: %w", err)
		}
	}
	return order, nil
}

func (s *Service) CreateOrder(ctx context.Context, order domain.Order) error {
	err := s.repo.Create(ctx, order)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateOrder):
			return fmt.Errorf("%w: %s", ErrOrderAlreadyExists, order.OrderUID)
		default:
			return fmt.Errorf("failed to create order: %w", err)
		}
	}
	return nil
}
