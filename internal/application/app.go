package application

import (
	"github.com/Killazius/L0/internal/application/kafka"
	"github.com/Killazius/L0/internal/config"
	"github.com/Killazius/L0/internal/repository/postgresql"
	"github.com/Killazius/L0/internal/service"
	"github.com/Killazius/L0/internal/transport/rest"
	"go.uber.org/zap"
)

type Application struct {
	logger   *zap.SugaredLogger
	Server   *rest.Server
	Consumer *kafka.Consumer
}

func New(logger *zap.SugaredLogger, cfg *config.Config) *Application {
	pool, err := postgresql.CreatePool(cfg.Postgres)
	if err != nil {
		logger.Fatalw("error creating postgres pool", "error", err)
	}
	orderRepo := postgresql.New(pool)
	orderService := service.New(orderRepo)
	return &Application{
		logger:   logger,
		Server:   rest.NewServer(logger, orderService, cfg.HTTPServer),
		Consumer: kafka.NewConsumer(logger, orderService, cfg.Kafka),
	}
}

func (app *Application) Run() {

}

func (app *Application) Stop() {}
