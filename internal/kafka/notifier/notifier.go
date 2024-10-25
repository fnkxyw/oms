package notifier

import (
	"github.com/IBM/sarama"
	"log"
	"time"
)

type Consumer struct{}

func (consumer *Consumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message received: key=%s value=%s", string(message.Key), string(message.Value))

		if err := processMessage(message); err != nil {
			log.Printf("Error processing message: %v", err)
			continue
		}

		session.MarkMessage(message, "")
		session.Commit()

	}
	return nil
}

func processMessage(msg *sarama.ConsumerMessage) error {
	log.Printf("[%s] Processing message: %s", time.Now().Format(time.RFC3339), string(msg.Value))
	return nil
}
