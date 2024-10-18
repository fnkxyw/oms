package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"log"
	"strconv"
	"time"
)

type AsyncProducer struct {
	producer sarama.AsyncProducer
	topic    string
}

func NewAsyncProducer(brokerList []string, topic string) (*AsyncProducer, error) {
	config := sarama.NewConfig()

	config.Net.MaxOpenRequests = 1
	config.Producer.Return.Successes = true
	config.Producer.Idempotent = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 10 * time.Millisecond
	config.Producer.Partitioner = sarama.NewHashPartitioner

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return &AsyncProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (k *AsyncProducer) SendMessage(order models.Order, event models.Event) error {
	model := Message{
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
		Key:       sarama.StringEncoder(strconv.Itoa(int(order.ID))),
		Value:     sarama.ByteEncoder(bytes),
		Timestamp: time.Now(),
	}

	k.producer.Input() <- msg

	go func() {
		for {
			select {
			case suc := <-k.producer.Successes():
				log.Printf("[KafkaProducer] Message sent successfully to partition: %d, offset: %d\n", suc.Partition, suc.Offset)

			case err := <-k.producer.Errors():
				log.Printf("Failed to send message: %v\n", err)
			}
		}
	}()

	return nil
}

func (k *AsyncProducer) SendMessages(orders []models.Order, event models.Event) error {
	doneChan := make(chan error, 1)
	defer close(doneChan)

	go func() {
		successCount := 0
		for {
			select {
			case suc := <-k.producer.Successes():
				log.Printf("[KafkaProducer] Message sent successfully to partition: %d, offset: %d\n", suc.Partition, suc.Offset)
				successCount++
				if successCount == len(orders) {
					doneChan <- nil
					return
				}
			case err := <-k.producer.Errors():
				log.Printf("Failed to send message: %v\n", err)
				doneChan <- err
				return
			}
		}
	}()

	for _, order := range orders {
		model := Message{
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
			Key:       sarama.StringEncoder(strconv.Itoa(int(order.ID))),
			Value:     sarama.ByteEncoder(bytes),
			Timestamp: time.Now(),
		}

		k.producer.Input() <- msg
	}

	return <-doneChan
}
