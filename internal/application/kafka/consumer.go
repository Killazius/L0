package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Killazius/L0/internal/config"
	"github.com/Killazius/L0/internal/domain"
	"github.com/Killazius/L0/internal/service"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"time"
)

type Consumer struct {
	reader  *kafka.Reader
	service *service.Service
	log     *zap.SugaredLogger
}

func NewConsumer(logger *zap.SugaredLogger, service *service.Service, cfg config.KafkaConfig) *Consumer {
	readerConfig := kafka.ReaderConfig{
		Brokers:         cfg.Brokers,
		Topic:           cfg.Topic,
		GroupID:         cfg.GroupID,
		SessionTimeout:  cfg.SessionTimeout,
		MaxWait:         cfg.MaxWait,
		MinBytes:        cfg.MinBytes,
		MaxBytes:        cfg.MaxBytes,
		RetentionTime:   7 * 24 * time.Hour,
		ReadLagInterval: 30 * time.Second,
		MaxAttempts:     cfg.MaxRetries,
		CommitInterval:  cfg.CommitInterval,
	}
	if cfg.AutoOffsetReset == "earliest" {
		readerConfig.StartOffset = kafka.FirstOffset
	} else {
		readerConfig.StartOffset = kafka.LastOffset
	}
	err := createTopicIfNotExists(cfg, 3, 1)
	if err != nil {
		logger.Fatal(err)
	}
	reader := kafka.NewReader(readerConfig)
	return &Consumer{
		reader:  reader,
		service: service,
		log:     logger,
	}
}

func (c *Consumer) Run(ctx context.Context) {
	c.log.Infow("starting kafka consumer",
		"topic", c.reader.Config().Topic,
		"group_id", c.reader.Config().GroupID,
	)
	for {
		select {
		case <-ctx.Done():
			c.log.Info("stopping kafka consumer due to context cancellation")
			return
		default:
			if err := c.Consume(ctx); err != nil {
				c.log.Errorw("Consume", "error", err)
			}
		}
	}
}

func (c *Consumer) Consume(ctx context.Context) error {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return err
	}
	var order domain.Order
	if err = json.Unmarshal(msg.Value, &order); err != nil {
		return err
	}
	c.log.Infow("read message", "order", order)
	if err = c.service.CreateOrder(ctx, order); err != nil {
		if errors.Is(err, service.ErrOrderAlreadyExists) {
			c.log.Warnw("order already exists", "order_uid", order.OrderUID)
			return nil
		}
		c.log.Errorw("failed to create order", "error", err)
		return err
	}
	if err = c.Commit(ctx, msg); err != nil {
		return err
	}
	return nil

}

func (c *Consumer) Commit(ctx context.Context, msg kafka.Message) error {
	if err := c.reader.CommitMessages(ctx, msg); err != nil {
		return err
	}
	return nil
}
func (c *Consumer) Stop() error {
	c.log.Info("stopping kafka consumer")
	return c.reader.Close()
}
