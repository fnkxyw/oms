package notifier

import (
	"context"
	"github.com/IBM/sarama"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/logger"
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
	ctx := context.Background()
	for message := range claim.Messages() {
		logger.Debugf(ctx, "Message received: key=%s value=%s", string(message.Key), string(message.Value))

		if err := processMessage(ctx, message); err != nil {
			logger.Errorf(ctx, "Error processing message: %v", err)
			continue
		}

		session.MarkMessage(message, "")
		session.Commit()

	}
	return nil
}

func processMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	logger.Infof(ctx, "[%s] Processing message: %s", time.Now().Format(time.RFC3339), string(msg.Value))
	return nil
}
