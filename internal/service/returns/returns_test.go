package returns

import (
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

	mockStorage := mocks.NewReturnStorageInterfaceMock(ctrl)

	mockStorage.GetReturnIDsMock.Expect().Return([]uint{55, 56, 57, 58, 59})
	mockStorage.GetReturnMock.When(uint(55)).Then(&models.Return{ID: 55, UserID: 1001}, true)
	mockStorage.GetReturnMock.When(uint(56)).Then(&models.Return{ID: 56, UserID: 1002}, true)
	mockStorage.GetReturnMock.When(uint(57)).Then(&models.Return{ID: 57, UserID: 1003}, true)
	mockStorage.GetReturnMock.When(uint(58)).Then(&models.Return{ID: 58, UserID: 1004}, true)
	mockStorage.GetReturnMock.When(uint(59)).Then(&models.Return{ID: 59, UserID: 1005}, true)

	err := ListReturns(mockStorage, 2, 2)
	assert.NoError(t, err)
}

func TestListReturns_InvalidPageNumber(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewReturnStorageInterfaceMock(ctrl)

	mockStorage.GetReturnIDsMock.Expect().Return([]uint{55, 56, 57, 58, 59})
	mockStorage.GetReturnMock.When(uint(55)).Then(&models.Return{ID: 55, UserID: 1001}, true)
	mockStorage.GetReturnMock.When(uint(56)).Then(&models.Return{ID: 56, UserID: 1002}, true)
	mockStorage.GetReturnMock.When(uint(57)).Then(&models.Return{ID: 57, UserID: 1003}, true)
	mockStorage.GetReturnMock.When(uint(58)).Then(&models.Return{ID: 58, UserID: 1004}, true)
	mockStorage.GetReturnMock.When(uint(59)).Then(&models.Return{ID: 59, UserID: 1005}, true)

	err := ListReturns(mockStorage, 2, -1)
	assert.Error(t, err)
}

func TestListReturns_NoMoreItems(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewReturnStorageInterfaceMock(ctrl)
	mockStorage.GetReturnIDsMock.Expect().Return([]uint{55, 56, 57, 58, 59})
	mockStorage.GetReturnMock.When(uint(55)).Then(&models.Return{ID: 55, UserID: 1001}, true)
	mockStorage.GetReturnMock.When(uint(56)).Then(&models.Return{ID: 56, UserID: 1002}, true)
	mockStorage.GetReturnMock.When(uint(57)).Then(&models.Return{ID: 57, UserID: 1003}, true)
	mockStorage.GetReturnMock.When(uint(58)).Then(&models.Return{ID: 58, UserID: 1004}, true)
	mockStorage.GetReturnMock.When(uint(59)).Then(&models.Return{ID: 59, UserID: 1005}, true)

	err := ListReturns(mockStorage, 2, 5) // запрашиваем несуществующую страницу
	assert.Error(t, err)
}

func TestListReturns_InvalidLimitNumber(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewReturnStorageInterfaceMock(ctrl)
	mockStorage.GetReturnIDsMock.Expect().Return([]uint{55, 56, 57, 58, 59})
	mockStorage.GetReturnMock.When(uint(55)).Then(&models.Return{ID: 55, UserID: 1001}, true)
	mockStorage.GetReturnMock.When(uint(56)).Then(&models.Return{ID: 56, UserID: 1002}, true)
	mockStorage.GetReturnMock.When(uint(57)).Then(&models.Return{ID: 57, UserID: 1003}, true)
	mockStorage.GetReturnMock.When(uint(58)).Then(&models.Return{ID: 58, UserID: 1004}, true)
	mockStorage.GetReturnMock.When(uint(59)).Then(&models.Return{ID: 59, UserID: 1005}, true)

	err := ListReturns(mockStorage, 0, 1) // неправильный лимит
	assert.Error(t, err)
}

func TestRefundOrder_OrderDoesNotExist(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockOrderStorage := mocks.NewOrderStorageInterfaceMock(ctrl)
	mockReturnStorage := mocks.NewReturnStorageInterfaceMock(ctrl)

	mockOrderStorage.GetOrderMock.When(uint(1)).Then(nil, false)

	err := RefundOrder(mockReturnStorage, mockOrderStorage, 1, 123)
	assert.Error(t, err)
}

func TestRefundOrder_OrderIsNotInPlaceState(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockOrderStorage := mocks.NewOrderStorageInterfaceMock(ctrl)
	mockReturnStorage := mocks.NewReturnStorageInterfaceMock(ctrl)

	mockOrderStorage.GetOrderMock.When(uint(101)).Then(&models.Order{
		ID:        101,
		State:     models.AcceptState,
		UserID:    123,
		PlaceDate: time.Now(),
	}, true)

	err := RefundOrder(mockReturnStorage, mockOrderStorage, 101, 123)
	assert.Error(t, err)
}

func TestRefundOrder_RefundTimeExpired(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockOrderStorage := mocks.NewOrderStorageInterfaceMock(ctrl)
	mockReturnStorage := mocks.NewReturnStorageInterfaceMock(ctrl)

	mockOrderStorage.GetOrderMock.When(uint(77)).Then(&models.Order{
		ID:        77,
		State:     models.PlaceState,
		UserID:    123,
		PlaceDate: time.Now().Add(-90 * time.Hour), // истёк срок
	}, true)

	err := RefundOrder(mockReturnStorage, mockOrderStorage, 77, 123)
	assert.Error(t, err)
}

func TestRefundOrder_IncorrectUserID(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockOrderStorage := mocks.NewOrderStorageInterfaceMock(ctrl)
	mockReturnStorage := mocks.NewReturnStorageInterfaceMock(ctrl)

	mockOrderStorage.GetOrderMock.When(uint(554)).Then(&models.Order{
		ID:        554,
		State:     models.PlaceState,
		UserID:    123,
		PlaceDate: time.Now(),
	}, true)

	err := RefundOrder(mockReturnStorage, mockOrderStorage, 554, 124) // неверный userId
	assert.Error(t, err)
}

func TestRefundOrder_SuccessfulRefund(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockOrderStorage := mocks.NewOrderStorageInterfaceMock(ctrl)
	mockReturnStorage := mocks.NewReturnStorageInterfaceMock(ctrl)

	mockOrderStorage.GetOrderMock.When(uint(88)).Then(&models.Order{
		ID:        88,
		State:     models.PlaceState,
		UserID:    123,
		PlaceDate: time.Now(),
	}, true)

	mockReturnStorage.AddReturnToStorageMock.When(&models.Return{
		ID:     88,
		UserID: 123,
	}).Then(nil)

	err := RefundOrder(mockReturnStorage, mockOrderStorage, 88, 123)
	assert.NoError(t, err)
}
