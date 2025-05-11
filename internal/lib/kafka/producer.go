package kafka

import (
	lib "bank-app-backend/internal/lib/logger"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type Producer struct {
	kafkaProducer *kafka.Producer
	topic         string
}

func NewProducer(brokers string, topic string) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
	if err != nil {
		return nil, err
	}

	return &Producer{kafkaProducer: p, topic: topic}, nil
}

func (p *Producer) SendEvent(key, value []byte) error {
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
	}

	err := p.kafkaProducer.Produce(msg, nil)
	if err != nil {
		lib.Log.Error("Produce error:", zap.Error(err))
	}
	return err
}
