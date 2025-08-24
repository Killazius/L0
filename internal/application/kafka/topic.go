package kafka

import (
	"github.com/Killazius/L0/internal/config"
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
)

func createTopicIfNotExists(cfg config.KafkaConfig, numPartitions, replicationFactor int) error {
	conn, err := kafka.Dial("tcp", cfg.Brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}

	defer controllerConn.Close()

	return controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             cfg.Topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	})
}
