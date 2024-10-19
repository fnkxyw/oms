package sync_producer

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/kafka"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"log"
	"strconv"
	"time"
)

const hubID int = 1010 // условно номер пвз куда отправляем извещения

type SyncProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewSyncProducer(brokerList []string, topic string) (*SyncProducer, error) {
	config := sarama.NewConfig()

	config.Net.MaxOpenRequests = 1
	config.Producer.Return.Successes = true
	config.Producer.Idempotent = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 10 * time.Millisecond
	config.Producer.Partitioner = sarama.NewHashPartitioner

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return &SyncProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (k *SyncProducer) SendMessage(ctx context.Context, order models.Order, event models.Event) error {
	model := kafka.Message{
		OrderID: order.ID,
		UserId:  order.UserID,
		Event:   event,
	}

	bytes, err := json.Marshal(model)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic:     k.topic,
		Key:       sarama.StringEncoder(strconv.Itoa(hubID)),
		Value:     sarama.ByteEncoder(bytes),
		Timestamp: time.Now(),
	}

	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message: %v\n", err)
		return err
	}
	log.Printf("[KafkaProducer] Message sent successfully to partition: %d, offset: %d\n", partition, offset)

	return nil
}

func (k *SyncProducer) SendMessages(ctx context.Context, orders []models.Order, event models.Event) error {
	for _, order := range orders {
		if err := k.SendMessage(ctx, order, event); err != nil {
			return err
		}
	}

	return nil
}

func (k *SyncProducer) Close() error {
	return k.producer.Close()
}
