package repository

import (
	"context"
	"github.com/Killazius/L0/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, order domain.Order) error
	Get(ctx context.Context, orderUID string) (*domain.Order, error)
}
