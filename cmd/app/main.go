package main

import (
	"l0/internal/application"
	"l0/internal/config"
	"l0/internal/logger"
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
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Info("Shutting down...")
}
