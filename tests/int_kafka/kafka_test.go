package integration_tests

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/kafka/notifier"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/kafka/sync_producer"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"log"
	"testing"
	"time"
)

func TestKafkaProducerConsumer(t *testing.T) {
	ctx := context.Background()

	kafkaContainer, err := kafka.Run(ctx,
		"confluentinc/cp-kafka:latest", kafka.WithClusterID("test-cluster"),
	)
	if err != nil {
		t.Fatalf("failed to start Kafka container: %v", err)
	}
	defer func() {
		if err := kafkaContainer.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	kafkaBrokers, err := kafkaContainer.Brokers(ctx)

	admin, err := sarama.NewClusterAdmin(kafkaBrokers, sarama.NewConfig())
	if err != nil {
		t.Fatalf("error Kafka admin create: %v", err)
	}
	defer admin.Close()

	topic := "test-topic"
	err = admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, false)
	if err != nil {
		t.Fatalf("error topic create %v", err)
	}

	producer, err := sync_producer.NewSyncProducer(kafkaBrokers, topic)
	if err != nil {
		t.Fatalf("error prooducer create: %v", err)
	}
	defer producer.Close()

	order := models.Order{
		ID:     1,
		UserID: 123,
	}
	event := models.AcceptEvent

	groupID := "test-group"
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumerGroup, err := sarama.NewConsumerGroup(kafkaBrokers, groupID, config)
	if err != nil {
		t.Fatalf("error consumer group create: %v", err)
	}

	consumer := &notifier.Consumer{}
	ready := make(chan struct{})

	go func() {
		close(ready)
		for {
			err := consumerGroup.Consume(ctx, []string{topic}, consumer)
			if err != nil {
				t.Fatalf("error message consume: %v", err)
			}
		}
	}()

	<-ready

	err = producer.SendMessage(ctx, order, event)
	if err != nil {
		t.Fatalf("error message sending: %v", err)
	}

	time.Sleep(5 * time.Second)
	assert.NoError(t, err)
}
