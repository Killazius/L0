package main

import (
	"context"
	"github.com/Killazius/L0/internal/application"
	"github.com/Killazius/L0/internal/config"
	"github.com/Killazius/L0/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log, err := logger.LoadFromConfig(cfg.Logger.Path)
	if err != nil {
		panic(err)
	}
	app := application.New(log, cfg)
	go app.Server.MustRun()
	go app.Consumer.Run(context.Background())
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Info("Shutting down...")
}
