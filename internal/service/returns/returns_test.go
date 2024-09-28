package returns

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/mocks"
	"testing"
	"time"
)

func TestListReturns_ValidReturnsWithPagination(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockStorage := mocks.NewStorageMock(ctrl)

	mockStorage.GetReturnsMock.Expect(ctx, models.RefundedState).Return([]*models.Order{
		{ID: 55, UserID: 1001, State: models.RefundedState},
		{ID: 56, UserID: 1002, State: models.RefundedState},
		{ID: 57, UserID: 1003, State: models.RefundedState},
		{ID: 58, UserID: 1004, State: models.RefundedState},
		{ID: 59, UserID: 1005, State: models.RefundedState},
	}, nil)

	err := ListReturns(ctx, mockStorage, 2, 2)
	assert.NoError(t, err)
}

func TestListReturns_InvalidPageNumber(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockStorage := mocks.NewStorageMock(ctrl)

	mockStorage.GetReturnsMock.Expect(ctx, models.RefundedState).Return([]*models.Order{
		{ID: 55, UserID: 1001, State: models.RefundedState},
		{ID: 56, UserID: 1002, State: models.RefundedState},
		{ID: 57, UserID: 1003, State: models.RefundedState},
		{ID: 58, UserID: 1004, State: models.RefundedState},
		{ID: 59, UserID: 1005, State: models.RefundedState},
	}, nil)

	err := ListReturns(ctx, mockStorage, 2, -1)
	assert.Error(t, err)
}

func TestListReturns_NoMoreItems(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockStorage := mocks.NewStorageMock(ctrl)

	mockStorage.GetReturnsMock.Expect(ctx, models.RefundedState).Return([]*models.Order{
		{ID: 55, UserID: 1001, State: models.RefundedState},
		{ID: 56, UserID: 1002, State: models.RefundedState},
		{ID: 57, UserID: 1003, State: models.RefundedState},
		{ID: 58, UserID: 1004, State: models.RefundedState},
		{ID: 59, UserID: 1005, State: models.RefundedState},
	}, nil)

	err := ListReturns(ctx, mockStorage, 2, 5)
	assert.Error(t, err)
}

func TestListReturns_InvalidLimitNumber(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockStorage := mocks.NewStorageMock(ctrl)

	mockStorage.GetReturnsMock.Expect(ctx, models.RefundedState).Return([]*models.Order{
		{ID: 55, UserID: 1001, State: models.RefundedState},
		{ID: 56, UserID: 1002, State: models.RefundedState},
		{ID: 57, UserID: 1003, State: models.RefundedState},
		{ID: 58, UserID: 1004, State: models.RefundedState},
		{ID: 59, UserID: 1005, State: models.RefundedState},
	}, nil)

	err := ListReturns(ctx, mockStorage, 0, 1) // неправильный лимит
	assert.Error(t, err)
}

func TestRefundOrder_OrderDoesNotExist(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockStorage := mocks.NewStorageMock(ctrl)

	mockStorage.GetItemMock.When(ctx, uint(1)).Then(&models.Order{}, false)

	err := RefundOrder(ctx, mockStorage, 1, 123)
	assert.Error(t, err)
}

func TestRefundOrder_OrderIsNotInPlaceState(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockStorage := mocks.NewStorageMock(ctrl)

	mockStorage.GetItemMock.When(ctx, uint(101)).Then(&models.Order{
		ID:        101,
		State:     models.AcceptState,
		UserID:    123,
		PlaceDate: time.Now(),
	}, true)

	err := RefundOrder(ctx, mockStorage, 101, 123)
	assert.Error(t, err)
}

func TestRefundOrder_RefundTimeExpired(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockStorage := mocks.NewStorageMock(ctrl)

	mockStorage.GetItemMock.When(ctx, uint(77)).Then(&models.Order{
		ID:        77,
		State:     models.PlaceState,
		UserID:    123,
		PlaceDate: time.Now().Add(-90 * time.Hour), // истёк срок
	}, true)

	err := RefundOrder(ctx, mockStorage, 77, 123)
	assert.Error(t, err)
}

func TestRefundOrder_IncorrectUserID(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockStorage := mocks.NewStorageMock(ctrl)

	mockStorage.GetItemMock.When(ctx, uint(554)).Then(&models.Order{
		ID:        554,
		State:     models.PlaceState,
		UserID:    123,
		PlaceDate: time.Now(),
	}, true)

	err := RefundOrder(ctx, mockStorage, 554, 124) // неверный userId
	assert.Error(t, err)
}

func TestRefundOrder_SuccessfulRefund(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockStorage := mocks.NewStorageMock(ctrl)

	mockStorage.GetItemMock.When(ctx, uint(88)).Then(&models.Order{
		ID:        88,
		State:     models.PlaceState,
		UserID:    123,
		PlaceDate: time.Now(),
	}, true)
	mockStorage.UpdateStateMock.Expect(ctx, uint(88), models.RefundedState).Return(nil)
	err := RefundOrder(ctx, mockStorage, 88, 123)
	assert.NoError(t, err)
}
