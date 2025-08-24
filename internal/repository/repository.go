package repository

import (
	"context"
	"errors"
	"github.com/Killazius/L0/internal/domain"
)

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrDeliveryNotFound = errors.New("delivery not found")
	ErrPaymentNotFound  = errors.New("payment not found")
	ErrItemsNotFound    = errors.New("items not found")
	ErrDuplicateOrder   = errors.New("duplicate order")
)

type Repository interface {
	Create(ctx context.Context, order domain.Order) error
	Get(ctx context.Context, orderUID string) (*domain.Order, error)
}
