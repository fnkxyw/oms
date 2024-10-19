package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/kafka/notifier"
	"log"
	"os"
	"os/signal"
)

const (
	topic         = "pvz.events-log"
	brokerAddress = "kafka0:29092"
)

func main() {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = false

	consumerGroup, err := sarama.NewConsumerGroup([]string{brokerAddress}, "pvz-consumer-group", config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer consumerGroup.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := &notifier.Consumer{}

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, consumer); err != nil {
				log.Fatalf("Error from consumer: %v", err)
			}
		}
	}()
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt)
	<-sigterm
	fmt.Println("Shutting down gracefully...")
}
