package orders

import (
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/mocks"
	"testing"
	"time"
)

func TestAcceptOrder_SuccessfulAcceptance(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	K := time.Now().Add(24 * time.Hour)
	T := time.Now().Unix()
	mockStorage := mocks.NewStorageMock(ctrl)

	mockStorage.IsConsistMock.When(uint(1)).Then(false)
	mockStorage.AddToStorageMock.Expect(&models.Order{
		ID:            1,
		KeepUntilDate: K,
		State:         models.AcceptState,
		AcceptTime:    T,
	})

	order := &models.Order{
		ID:            1,
		KeepUntilDate: K,
		State:         models.AcceptState,
		AcceptTime:    T,
	}

	err := AcceptOrder(mockStorage, order)
	assert.NoError(t, err)
}

func TestAcceptOrder_OrderAlreadyExists(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)

	mockStorage.IsConsistMock.When(uint(4)).Then(true)

	order := &models.Order{
		ID:            4,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
		State:         models.AcceptState,
	}

	err := AcceptOrder(mockStorage, order)
	assert.Error(t, err)
}

func TestAcceptOrder_OrderExpired(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)

	mockStorage.IsConsistMock.When(uint(2)).Then(false)

	order := &models.Order{
		ID:            2,
		KeepUntilDate: time.Now().Add(-24 * time.Hour),
		State:         models.AcceptState,
	}

	err := AcceptOrder(mockStorage, order)
	assert.Error(t, err)
}

func TestPlaceOrder_EmptyIDs(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)

	err := PlaceOrder(mockStorage, []uint{})
	assert.Error(t, err)
}

func TestPlaceOrder_NoConsist(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)
	mockStorage.GetItemMock.Expect(uint(1)).Return(&models.Order{
		ID:            1,
		State:         models.NewState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}, false)

	err := PlaceOrder(mockStorage, []uint{1})
	assert.Error(t, err)
}

func TestPlaceOrder_Success(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)
	mockStorage.GetItemMock.When(uint(5)).Then(&models.Order{
		ID:            5,
		State:         models.AcceptState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}, true)

	err := PlaceOrder(mockStorage, []uint{5})
	assert.NoError(t, err)
}

func TestReturnOrder_SuccessfulRefundedState(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)
	mockStorage.GetItemMock.When(uint(199)).Then(&models.Order{
		ID:     199,
		UserID: 1,
		State:  models.RefundedState,
	}, true)

	err := ReturnOrder(mockStorage, 199)
	assert.NoError(t, err)
}

// Тест возвращения заказа с истекшей датой
func TestReturnOrder_SuccessfulReturnExpired(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)
	mockStorage.GetItemMock.When(uint(750)).Then(&models.Order{
		ID:            750,
		UserID:        10,
		State:         models.AcceptState,
		KeepUntilDate: time.Now().Add(-47 * time.Hour),
	}, true)

	err := ReturnOrder(mockStorage, 750)
	assert.NoError(t, err)
}

// Тест, когда заказ не найден
func TestReturnOrder_NoOrderFound(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)
	mockStorage.GetItemMock.When(uint(191)).Then(&models.Order{}, false)

	err := ReturnOrder(mockStorage, 191)
	assert.Error(t, err)
}

// Тест, когда заказ не может быть возвращен
func TestReturnOrder_OrderCannotBeReturned(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)
	mockStorage.GetItemMock.When(uint(150)).Then(&models.Order{
		ID:            150,
		UserID:        1,
		State:         models.AcceptState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}, true)

	err := ReturnOrder(mockStorage, 150)
	assert.Error(t, err)
}

// Тест, когда заказов нет в наличии
func TestListOrders_NoOrdersAvailable(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)

	mockStorage.GetIDsMock.Expect().Return([]uint{})

	err := ListOrders(mockStorage, uint(912), 1, false)

	assert.NoError(t, err)
}

func TestListOrders_NoConsist(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewStorageMock(ctrl)
	mockStorage.GetIDsMock.Expect().Return([]uint{912})
	mockStorage.GetItemMock.When(uint(912)).Then(&models.Order{}, false)

	err := ListOrders(mockStorage, uint(1), 1, false)
	assert.NoError(t, err)
}
