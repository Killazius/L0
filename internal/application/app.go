package application

import (
	"go.uber.org/zap"
	"l0/internal/config"
	"l0/internal/repository/postgresql"
	"l0/internal/service"
	"l0/internal/transport"
)

type Application struct {
	logger *zap.SugaredLogger
	Server *transport.Server
}

func New(logger *zap.SugaredLogger, cfg *config.Config) *Application {
	pool, err := config.CreatePool(cfg.Postgres)
	if err != nil {
		logger.Fatalw("error creating postgres pool", "error", err)
	}
	orderRepo := postgresql.New(pool)
	orderService := service.New(orderRepo)
	httpServer := transport.NewServer(logger, orderService, cfg.HTTPServer)
	return &Application{
		logger: logger,
		Server: httpServer,
	}
}
