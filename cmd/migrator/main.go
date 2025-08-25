package main

import (
	"errors"
	"flag"
	"github.com/Killazius/L0/internal/config"
	"github.com/Killazius/L0/internal/logger"
	"github.com/Killazius/L0/internal/repository/postgresql"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
	"os"
)

var command = flag.String("command", "up", "goose command (up, down, status, etc.)")

func main() {
	cfg := config.MustLoad()
	log, err := logger.LoadFromConfig(cfg.Logger.Path)
	if err != nil {
		if errors.Is(err, logger.ErrDefaultLogger) {
			log.Warnw("using default logger because config file not found",
				"config_path", cfg.Logger.Path)
		} else {
			log.Fatal(err)
		}
	}
	if _, err = os.Stat(cfg.Postgres.MigrationsPath); os.IsNotExist(err) {
		log.Fatalw("migrations directory does not exist", "path", cfg.Postgres.MigrationsPath)
	}
	if *command == "" {
		if len(flag.Args()) > 0 {
			*command = flag.Args()[0]
		} else {
			log.Info("no goose command provided. Usage: -command <command> or provide command as argument")
		}
	}

	pool, err := postgresql.CreatePool(cfg.Postgres)
	if err != nil {
		log.Fatalw("error creating postgres pool", "error", err)
	}
	defer pool.Close()
	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()
	if err := goose.Run(*command, sqlDB, cfg.Postgres.MigrationsPath); err != nil {
		log.Fatalw("failed to run goose command", "command", *command, "error", err)
	}

	log.Infow("goose command executed successfully", "command", *command)
}
