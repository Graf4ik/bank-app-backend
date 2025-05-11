package kafka

import (
	lib "bank-app-backend/internal/lib/logger"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

func RunConsumer(brokers, topic, groupID string) {
	c, _ := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	err := c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return
	}

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			lib.Log.Info("Received message", zap.String("value", string(msg.Value)))
		}
	}
}
