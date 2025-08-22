package repository

import (
	"context"
	"l0/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, order domain.Order) error
}
