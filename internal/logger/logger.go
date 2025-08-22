package logger

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
)

func LoadFromConfig(path string) (*zap.SugaredLogger, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("logger config file does not exist: %s", path)
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
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to build logger from config %q: %w", path, err)
	}
	undo := zap.ReplaceGlobals(logger)
	defer undo()
	return logger.Sugar(), nil
}
