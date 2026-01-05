package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"paste/internal/model"
	"paste/internal/service"
	"time"

	"github.com/google/uuid"
	kgo "github.com/segmentio/kafka-go"
)

type KafkaGoConsumer struct {
	r           *kgo.Reader
	scanService *service.ScanService
}

func NewKafkaGoConsumer(brokers []string, topic, groupID string, scanService *service.ScanService) *KafkaGoConsumer {
	r := kgo.NewReader(kgo.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
		StartOffset:    kgo.FirstOffset,
	})

	log.Printf("Kafka consumer создан: brokers=%v, topic=%s, groupID=%s", brokers, topic, groupID)

	return &KafkaGoConsumer{
		r:           r,
		scanService: scanService,
	}
}

func (c *KafkaGoConsumer) Start(ctx context.Context) error {
	log.Println("Kafka consumer начал читать сообщения...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer: получен сигнал остановки")
			return ctx.Err()
		default:
			msg, err := c.r.FetchMessage(ctx)
			if err != nil {
				if err == context.Canceled || err == context.DeadlineExceeded {
					return err
				}
				log.Printf("Ошибка чтения из Kafka: %v. Повторная попытка через 5 секунд...", err)
				time.Sleep(5 * time.Second)
				continue
			}

			log.Printf("Получено сообщение из Kafka: partition=%d, offset=%d, key=%s",
				msg.Partition, msg.Offset, string(msg.Key))

			if err := c.processMessage(ctx, msg); err != nil {
				log.Printf("Ошибка обработки сообщения: %v", err)
			} else {
				// Коммитим сообщение только после успешной обработки
				if err := c.r.CommitMessages(ctx, msg); err != nil {
					log.Printf("Ошибка коммита: %v", err)
				}
			}
		}
	}
}

func (c *KafkaGoConsumer) processMessage(ctx context.Context, msg kgo.Message) error {
	var scan model.Scan
	if err := json.Unmarshal(msg.Value, &scan); err != nil {
		return fmt.Errorf("не удалось распарсить JSON: %w", err)
	}

	scan.ID = uuid.New()
	scan.CreatedAt = time.Now()

	if err := c.scanService.CreateScan(ctx, &scan); err != nil {
		return fmt.Errorf("не удалось сохранить scan в БД: %w", err)
	}

	log.Printf("Обработано сообщение, ключ: %s", scan.FullName)
	return nil
}

func (c *KafkaGoConsumer) Close() error {
	return c.r.Close()
}
