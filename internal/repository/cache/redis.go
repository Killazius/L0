package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Killazius/L0/internal/config"
	"github.com/Killazius/L0/internal/domain"
	"github.com/Killazius/L0/internal/repository"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"time"
)

type Cache struct {
	client *redis.Client
}

const defaultTTL = time.Hour * 24

func CreateClient(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		DB:       cfg.DB,
		Password: cfg.Password,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}
	return client, nil
}
func New(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) Set(ctx context.Context, order *domain.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	err = c.client.Set(ctx, order.OrderUID, data, defaultTTL).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) Get(ctx context.Context, orderUID string) (*domain.Order, error) {

	orderJSON, err := c.client.Get(ctx, orderUID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, repository.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order from cache: %w", err)
	}

	var order domain.Order
	err = json.Unmarshal([]byte(orderJSON), &order)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	return &order, nil
}

func (c *Cache) Restore(ctx context.Context, repo repository.OrderRepository, workers int) error {
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
