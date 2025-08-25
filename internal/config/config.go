package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
	"time"
)

type Config struct {
	Postgres   PostgresConfig `yaml:"postgres"`
	Logger     LoggerConfig   `yaml:"logger"`
	HTTPServer HTTPConfig     `yaml:"http_server"`
	Kafka      KafkaConfig    `yaml:"kafka"`
}

type HTTPConfig struct {
	Port        string        `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	Host        string        `yaml:"host" env:"HTTP_HOST" env-default:"localhost"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT" env-default:"30s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
}

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	Username string `env:"POSTGRES_USER" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	Database string `env:"POSTGRES_DB" env-default:"postgres"`
	SSLMode  string `env:"POSTGRES_SSL_MODE" env-default:"disable"`

	MaxOpenConns    int32         `yaml:"max_open_conns" env-default:"25"`
	MaxIdleConns    int32         `yaml:"max_idle_conns" env-default:"5"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env-default:"1h"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" env-default:"30m"`
	Timeout         time.Duration `yaml:"timeout" env-default:"5s"`
	MigrationsPath  string        `yaml:"migrations_path" env-default:"./migrations"`
}

type LoggerConfig struct {
	Path string `yaml:"path"`
}

type KafkaConfig struct {
	Brokers          []string      `yaml:"brokers"  env-default:"kafka:9092" env-separator:","`
	Topic            string        `yaml:"topic" env-default:"orders"`
	GroupID          string        `yaml:"group_id"  env-default:"order-service-group"`
	AutoOffsetReset  string        `yaml:"auto_offset_reset" env-default:"earliest"`
	SessionTimeout   time.Duration `yaml:"session_timeout" env-default:"30s"`
	MaxWait          time.Duration `yaml:"max_wait" env-default:"10s"`
	MinBytes         int           `yaml:"min_bytes"  env-default:"10240"`
	MaxBytes         int           `yaml:"max_bytes" env-default:"10485760"`
	MaxRetries       int           `yaml:"max_retries"  env-default:"3"`
	RetryBackoff     time.Duration `yaml:"retry_backoff" env-default:"100ms"`
	EnableAutoCommit bool          `yaml:"enable_auto_commit" env-default:"false"`
	CommitInterval   time.Duration `yaml:"commit_interval" env-default:"1s"`
}

const (
	defaultConfigPath = "config/config.yaml"
)

func MustLoad() *Config {
	cfg, err := load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return cfg
}

func load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	configPath := getConfigPath()
	if configPath != "" {
		fileInfo, err := os.Stat(configPath)
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file %s does not exist", configPath)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to access config file: %w", err)
		}
		if fileInfo.IsDir() {
			return nil, fmt.Errorf("config path is a directory, not a file: %s", configPath)
		}
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read env vars: %w", err)
	}

	return &cfg, nil
}

func getConfigPath() string {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = defaultConfigPath
	}
	return configPath
}

func (h *HTTPConfig) GetAddr() string {
	return net.JoinHostPort(h.Host, h.Port)
}
func (p *PostgresConfig) GetURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", p.Username, p.Password, p.Host, p.Port, p.Database, p.SSLMode)
}
