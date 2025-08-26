package main

import (
	"context"
	"errors"
	"github.com/Killazius/L0/internal/application"
	"github.com/Killazius/L0/internal/config"
	"github.com/Killazius/L0/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := config.MustLoad()
	log, err := logger.LoadFromConfig(cfg.Logger.Path)
	if err != nil {
		if errors.Is(err, logger.ErrDefaultLogger) {
			log.Warnw("using default logger because config file not found",
				"config_path", cfg.Logger.Path)
		} else {
			panic(err)
		}
	}
	app := application.New(log, cfg)
	app.Run(ctx)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	cancel()
	app.Stop()
}
