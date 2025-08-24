package application

import (
	"github.com/Killazius/L0/internal/application/kafka"
	"github.com/Killazius/L0/internal/config"
	"github.com/Killazius/L0/internal/repository/postgresql"
	"github.com/Killazius/L0/internal/service"
	"github.com/Killazius/L0/internal/transport"
	"go.uber.org/zap"
)

type Application struct {
	logger   *zap.SugaredLogger
	Server   *transport.Server
	Consumer *kafka.Consumer
}

func New(logger *zap.SugaredLogger, cfg *config.Config) *Application {
	pool, err := config.CreatePool(cfg.Postgres)
	if err != nil {
		logger.Fatalw("error creating postgres pool", "error", err)
	}
	orderRepo := postgresql.New(pool)
	orderService := service.New(orderRepo)
	httpServer := transport.NewServer(logger, orderService, cfg.HTTPServer)
	consumer := kafka.NewConsumer(logger, orderService, cfg.Kafka)
	return &Application{
		logger:   logger,
		Server:   httpServer,
		Consumer: consumer,
	}
}

func (app *Application) Run() {

}

func (app *Application) Stop() {}
