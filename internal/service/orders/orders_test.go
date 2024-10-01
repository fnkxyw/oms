package orders

import (
	"context"
	"fmt"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/mocks"

	"testing"
	"time"
)

func TestAcceptOrder(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	defer mc.Finish()

	mockStorage := mocks.NewFacadeMock(mc)
	ctx := context.TODO()
	order := &models.Order{ID: 1}

	mockStorage.AcceptOrderMock.Expect(ctx, order).Return(nil)

	err := AcceptOrder(ctx, mockStorage, order)

	require.NoError(t, err)
	mockStorage.MinimockFinish()
}

func TestAcceptOrder_OrderAlreadyExists(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewFacadeMock(ctrl)

	order := &models.Order{
		ID:            4,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
		State:         models.AcceptState,
	}

	mockStorage.AcceptOrderMock.Expect(ctx, order).Return(e.ErrIsConsist)

	err := AcceptOrder(ctx, mockStorage, order)
	assert.Error(t, err)
}

func TestAcceptOrder_OrderExpired(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewFacadeMock(ctrl)

	order := &models.Order{
		ID:            2,
		KeepUntilDate: time.Now().Add(-24 * time.Hour),
		State:         models.AcceptState,
	}

	mockStorage.AcceptOrderMock.Expect(ctx, order).Return(e.ErrTimeExpired)

	err := AcceptOrder(ctx, mockStorage, order)
	assert.Error(t, err)
}

func TestPlaceOrder_EmptyIDs(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewFacadeMock(ctrl)

	err := PlaceOrder(ctx, mockStorage, []uint{})
	assert.Error(t, err)
}

func TestPlaceOrder_NoConsist(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewFacadeMock(ctrl)
	mockStorage.PlaceOrderMock.Expect(ctx, []uint{1}).Return(e.ErrNoConsist)

	err := PlaceOrder(ctx, mockStorage, []uint{1})
	assert.Error(t, err)
}

func TestPlaceOrder_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()
	mockStorage := mocks.NewFacadeMock(ctrl)

	mockStorage.PlaceOrderMock.Expect(ctx, []uint{5}).Return(nil)

	err := PlaceOrder(ctx, mockStorage, []uint{5})

	assert.NoError(t, err)
}

func TestReturnOrder_SuccessfulRefundedState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewFacadeMock(ctrl)
	mockStorage.ReturnOrderMock.Expect(ctx, uint(199)).Return(nil)

	err := ReturnOrder(ctx, mockStorage, 199)
	assert.NoError(t, err)
}

func TestReturnOrder_SuccessfulReturnExpired(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewFacadeMock(ctrl)
	mockStorage.ReturnOrderMock.Expect(ctx, uint(750)).Return(nil)

	err := ReturnOrder(ctx, mockStorage, 750)
	assert.NoError(t, err)
}

func TestReturnOrder_NoOrderFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewFacadeMock(ctrl)
	mockStorage.ReturnOrderMock.Expect(ctx, uint(191)).Return(fmt.Errorf("order not found"))

	err := ReturnOrder(ctx, mockStorage, 191)
	assert.Error(t, err)
}

func TestReturnOrder_OrderCannotBeReturned(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewFacadeMock(ctrl)
	mockStorage.ReturnOrderMock.Expect(ctx, uint(150)).Return(fmt.Errorf("order cannot be returned"))

	err := ReturnOrder(ctx, mockStorage, 150)
	assert.Error(t, err)
}

// Тест, когда заказов нет в наличии
func TestListOrders_NoOrdersAvailable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewFacadeMock(ctrl)

	mockStorage.ListOrdersMock.Expect(ctx, uint(912), false).Return(nil, nil)

	err := ListOrders(ctx, mockStorage, uint(912), 1, false)

	assert.NoError(t, err)
}

func TestListOrders_NoConsist(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewFacadeMock(ctrl)

	mockStorage.ListOrdersMock.Expect(ctx, uint(1), false).Return([]models.Order{
		{ID: 912, UserID: 1, State: models.AcceptState},
	}, nil)

	err := ListOrders(ctx, mockStorage, uint(1), 1, false)
	assert.NoError(t, err)
}
