package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
)

var ErrDefaultLogger = errors.New("default loggers")

func LoadFromConfig(path string) (*zap.SugaredLogger, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		defaultLogger, _ := zap.NewDevelopment()
		sugar := defaultLogger.Sugar()
		return sugar, ErrDefaultLogger
	}
	configData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg zap.Config
	if err = json.Unmarshal(configData, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger from config %q: %w", path, err)
	}
	undo := zap.ReplaceGlobals(logger)
	defer undo()
	return logger.Sugar(), nil
}
