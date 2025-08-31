package service

import (
	"context"
	"errors"
	"github.com/Killazius/L0/internal/domain"
	"sync"
	"time"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrInvalidOrderData   = errors.New("invalid order data")
)

const (
	cacheTimeout = 3 * time.Second
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	Get(ctx context.Context, orderUID string) (*domain.Order, error)
	GetAll(ctx context.Context) ([]domain.Order, error)
}

type OrderCache interface {
	Set(ctx context.Context, order *domain.Order) error
	Get(ctx context.Context, orderUID string) (*domain.Order, error)
}

type Service struct {
	repo  OrderRepository
	cache OrderCache
	wg    sync.WaitGroup
}

func New(repo OrderRepository, cache OrderCache) *Service {
	return &Service{repo: repo, cache: cache}
}
