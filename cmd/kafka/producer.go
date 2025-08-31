package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/Killazius/L0/internal/lib/test"
	"github.com/Killazius/L0/internal/logger"
	"github.com/segmentio/kafka-go"
	"time"
)

var messageCount = flag.Int("m", 1, "count of messages to send")

const defaultLoggerPath = "config/logger.json"

func main() {
	flag.Parse()
	brokers := []string{"localhost:29092"}
	topic := "orders"
	producer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer producer.Close()
	log, err := logger.LoadFromConfig(defaultLoggerPath)
	if err != nil {
		panic(err)
	}
	log.Infow("start producer", "topic", topic, "messageCount", *messageCount)

	for range *messageCount {
		order := test.GenerateOrder()

		data, err := json.Marshal(order)
		if err != nil {
			log.Errorw("failed to marshal order", "error", err, "order", order)
			continue
		}

		err = producer.WriteMessages(context.Background(),
			kafka.Message{
				Value: data,
				Key:   []byte(order.OrderUID),
			},
		)

		if err != nil {
			log.Errorw("failed to write messages", "error", err, "order", order)
		} else {
			log.Infow("success", "order", order)
		}

		time.Sleep(500 * time.Millisecond)
	}
	log.Info("producer stopped")

}
