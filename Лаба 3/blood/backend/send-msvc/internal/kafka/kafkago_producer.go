package kafka

import (
	"context"
	"log"
	"time"

	kgo "github.com/segmentio/kafka-go"
)

type KafkaGoProducer struct {
	w *kgo.Writer
}

func NewKafkaGoProducer(brokers []string, topic string) *KafkaGoProducer {
	return &KafkaGoProducer{
		w: &kgo.Writer{
			Addr:         kgo.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kgo.LeastBytes{},
			BatchTimeout: 50 * time.Millisecond,
			RequiredAcks: kgo.RequireOne,
		},
	}
}

func (p *KafkaGoProducer) Publish(ctx context.Context, key string, value []byte) error {
	err := p.w.WriteMessages(ctx, kgo.Message{
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	})
	if err != nil {
		return err
	}
	log.Println("Опубликовано сообщение в Kafka, ключ:", key)
	return nil
}

func (p *KafkaGoProducer) Close() error {
	return p.w.Close()
}
