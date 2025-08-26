package application

import (
	"context"
	"github.com/Killazius/L0/internal/application/kafka"
	"github.com/Killazius/L0/internal/config"
	"github.com/Killazius/L0/internal/repository/cache"
	"github.com/Killazius/L0/internal/repository/postgresql"
	"github.com/Killazius/L0/internal/service"
	"github.com/Killazius/L0/internal/transport/rest"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Application struct {
	log         *zap.SugaredLogger
	Server      *rest.Server
	Consumer    *kafka.Consumer
	pool        *pgxpool.Pool
	cacheClient *redis.Client
	wg          sync.WaitGroup
}

func New(log *zap.SugaredLogger, cfg *config.Config) *Application {
	pool, err := postgresql.CreatePool(cfg.Postgres)
	if err != nil {
		log.Fatalw("error creating postgres pool", "error", err)
	}
	client, err := cache.CreateClient(cfg.Redis)
	if err != nil {
		log.Fatalw("error creating redis client", "error", err)
	}
	orderRepo := postgresql.New(pool)
	orderCache := cache.New(client)

	if err = orderCache.Restore(context.Background(), orderRepo, 10); err != nil {
		log.Fatalw("error restoring order", "error", err)
	}

	orderService := service.New(orderRepo, orderCache)
	return &Application{
		log:         log,
		Server:      rest.NewServer(log, orderService, cfg.HTTPServer),
		Consumer:    kafka.NewConsumer(log, orderService, cfg.Kafka),
		pool:        pool,
		cacheClient: client,
	}
}

func (a *Application) Run(ctx context.Context) {
	a.wg.Go(a.Server.MustRun)
	a.wg.Go(func() {
		a.Consumer.Run(ctx)
	})
}

func (a *Application) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	a.log.Info("closing HTTP server")
	if err := a.Server.Close(ctx); err != nil {
		a.log.Errorw("failed to stop HTTP server gracefully", "error", err)
	}

	a.log.Info("closing Kafka consumer")
	if err := a.Consumer.Close(); err != nil {
		a.log.Errorw("failed to stop Kafka consumer gracefully", "error", err)
	}
	a.log.Info("closing database connections")
	a.pool.Close()
	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-ctx.Done():
		a.log.Errorw("failed to stop gracefully", "error", ctx.Err())
	}
}
