package sync_producer

import (
	"context"
	"fmt"
	"github.com/IBM/sarama/mocks"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"testing"
)

func TestSendMessage_Success(t *testing.T) {
	topic := "test_topic"
	order := models.Order{ID: 1, UserID: 2}
	event := models.AcceptEvent

	mockProducer := mocks.NewSyncProducer(t, nil)
	defer mockProducer.Close()

	syncProducer := &SyncProducer{
		producer: mockProducer,
		topic:    topic,
	}

	mockProducer.ExpectSendMessageAndSucceed()

	ctx := context.Background()

	err := syncProducer.SendMessage(ctx, order, event)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

}

func TestSendMessage_Failure(t *testing.T) {
	topic := "test_topic"
	order := models.Order{ID: 1, UserID: 2}
	event := models.AcceptEvent

	mockProducer := mocks.NewSyncProducer(t, nil)
	defer mockProducer.Close()

	syncProducer := &SyncProducer{
		producer: mockProducer,
		topic:    topic,
	}

	mockProducer.ExpectSendMessageAndFail(fmt.Errorf("failed to send message"))

	ctx := context.Background()

	err := syncProducer.SendMessage(ctx, order, event)

	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}
