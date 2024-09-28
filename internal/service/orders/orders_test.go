package orders

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/mocks"

	"testing"
	"time"
)

func TestAcceptOrder_SuccessfulAcceptance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	K := time.Now().Add(24 * time.Hour)
	T := time.Now().Unix()
	mockStorage := mocks.NewStorageMock(ctrl)

	mockStorage.IsConsistMock.When(ctx, uint(1)).Then(false)
	mockStorage.AddToStorageMock.Expect(ctx, &models.Order{
		ID:            1,
		KeepUntilDate: K,
		State:         models.AcceptState,
		AcceptTime:    T,
	}).Return(nil)

	order := &models.Order{
		ID:            1,
		KeepUntilDate: K,
		State:         models.AcceptState,
		AcceptTime:    T,
	}

	err := AcceptOrder(ctx, mockStorage, order)
	assert.NoError(t, err)
}

func TestAcceptOrder_OrderAlreadyExists(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)

	mockStorage.IsConsistMock.When(ctx, uint(4)).Then(true)

	order := &models.Order{
		ID:            4,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
		State:         models.AcceptState,
	}

	err := AcceptOrder(ctx, mockStorage, order)
	assert.Error(t, err)
}

func TestAcceptOrder_OrderExpired(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)

	mockStorage.IsConsistMock.When(ctx, uint(2)).Then(false)

	order := &models.Order{
		ID:            2,
		KeepUntilDate: time.Now().Add(-24 * time.Hour),
		State:         models.AcceptState,
	}

	err := AcceptOrder(ctx, mockStorage, order)
	assert.Error(t, err)
}

func TestPlaceOrder_EmptyIDs(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)

	err := PlaceOrder(ctx, mockStorage, []uint{})
	assert.Error(t, err)
}

func TestPlaceOrder_NoConsist(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)
	mockStorage.GetItemMock.Expect(ctx, uint(1)).Return(&models.Order{
		ID:            1,
		State:         models.NewState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}, false)

	err := PlaceOrder(ctx, mockStorage, []uint{1})
	assert.Error(t, err)
}

func TestPlaceOrder_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()
	currentTime := time.Now()
	mockStorage := mocks.NewStorageMock(ctrl)
	mockStorage.GetItemMock.When(ctx, uint(5)).Then(&models.Order{
		ID:            5,
		State:         models.AcceptState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}, true)

	mockStorage.UpdateBeforePlaceMock.Set(func(ctx context.Context, id uint, state models.State, placeTime time.Time) error {
		assert.Equal(t, uint(5), id)
		assert.Equal(t, models.PlaceState, state)
		assert.WithinDuration(t, currentTime, placeTime, time.Second)

		return nil
	})

	err := PlaceOrder(ctx, mockStorage, []uint{5})

	assert.NoError(t, err)
}

func TestReturnOrder_SuccessfulRefundedState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)
	mockStorage.GetItemMock.When(ctx, uint(199)).Then(&models.Order{
		ID:     199,
		UserID: 1,
		State:  models.RefundedState,
	}, true)
	mockStorage.DeleteFromStorageMock.Expect(ctx, uint(199)).Return(nil)

	err := ReturnOrder(ctx, mockStorage, 199)
	assert.NoError(t, err)
}

func TestReturnOrder_SuccessfulReturnExpired(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)
	mockStorage.GetItemMock.When(ctx, uint(750)).Then(&models.Order{
		ID:            750,
		UserID:        10,
		State:         models.AcceptState,
		KeepUntilDate: time.Now().Add(-47 * time.Hour),
	}, true)

	mockStorage.DeleteFromStorageMock.Expect(ctx, uint(750)).Return(nil)

	err := ReturnOrder(ctx, mockStorage, 750)
	assert.NoError(t, err)
}

// Тест, когда заказ не найден
func TestReturnOrder_NoOrderFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)
	mockStorage.GetItemMock.When(ctx, uint(191)).Then(&models.Order{}, false)

	err := ReturnOrder(ctx, mockStorage, 191)
	assert.Error(t, err)
}

// Тест, когда заказ не может быть возвращен
func TestReturnOrder_OrderCannotBeReturned(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)
	mockStorage.GetItemMock.When(ctx, uint(150)).Then(&models.Order{
		ID:            150,
		UserID:        1,
		State:         models.AcceptState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}, true)

	err := ReturnOrder(ctx, mockStorage, 150)
	assert.Error(t, err)
}

// Тест, когда заказов нет в наличии
func TestListOrders_NoOrdersAvailable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)

	mockStorage.GetByUserIdMock.Expect(ctx, uint(912)).Return(nil, nil)

	err := ListOrders(ctx, mockStorage, uint(912), 1, false)

	assert.NoError(t, err)
}

func TestListOrders_NoConsist(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)

	mockStorage.GetByUserIdMock.Expect(ctx, uint(1)).Return([]*models.Order{
		{ID: 912, UserID: 1, State: models.AcceptState},
	}, nil)

	err := ListOrders(ctx, mockStorage, uint(1), 1, false)
	assert.NoError(t, err)
}
