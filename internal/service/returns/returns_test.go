package returns

import (
	"context"
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/mocks"
	"testing"
)

func TestListReturns_ValidReturnsWithPagination(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockStorage := mocks.NewFacadeMock(ctrl)
	mockStorage.ListReturnsMock.Expect(ctx, 2, 2).Return([]models.Order{
		{ID: 55, UserID: 1001, State: models.RefundedState},
		{ID: 56, UserID: 1002, State: models.RefundedState},
	}, nil)

	err := ListReturns(ctx, mockStorage, 2, 2)
	assert.NoError(t, err)

}

func TestListReturns_InvalidPageNumber(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockStorage := mocks.NewFacadeMock(ctrl)

	mockStorage.ListReturnsMock.Expect(ctx, 2, -1).Return(nil, errors.New("invalid page number"))

	err := ListReturns(ctx, mockStorage, 2, -1)
	assert.Error(t, err)
}

func TestListReturns_NoMoreItems(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockStorage := mocks.NewFacadeMock(ctrl)

	mockStorage.ListReturnsMock.Expect(ctx, 2, 5).Return([]models.Order{}, nil)

	err := ListReturns(ctx, mockStorage, 2, 5)
	assert.NoError(t, err)
}

func TestListReturns_InvalidLimitNumber(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockStorage := mocks.NewFacadeMock(ctrl)

	mockStorage.ListReturnsMock.Expect(ctx, 0, 1).Return(nil, errors.New("invalid limit number"))

	err := ListReturns(ctx, mockStorage, 0, 1)
	assert.Error(t, err)
}

func TestRefundOrder_OrderDoesNotExist(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockStorage := mocks.NewFacadeMock(ctrl)

	mockStorage.RefundOrderMock.Expect(ctx, uint(1), uint(123)).Return(e.ErrNoConsist)

	err := RefundOrder(ctx, mockStorage, 1, 123)
	assert.ErrorIs(t, err, e.ErrNoConsist)
}

func TestRefundOrder_OrderIsNotInPlaceState(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockStorage := mocks.NewFacadeMock(ctrl)

	mockStorage.RefundOrderMock.Expect(ctx, uint(101), uint(123)).Return(e.ErrNotPlace)

	err := RefundOrder(ctx, mockStorage, 101, 123)
	assert.ErrorIs(t, err, e.ErrNotPlace)
}

func TestRefundOrder_RefundTimeExpired(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockStorage := mocks.NewFacadeMock(ctrl)

	mockStorage.RefundOrderMock.Expect(ctx, uint(77), uint(123)).Return(e.ErrTimeExpired)

	err := RefundOrder(ctx, mockStorage, 77, 123)
	assert.ErrorIs(t, err, e.ErrTimeExpired)
}

func TestRefundOrder_IncorrectUserID(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockStorage := mocks.NewFacadeMock(ctrl)

	mockStorage.RefundOrderMock.Expect(ctx, uint(554), uint(124)).Return(e.ErrIncorrectUserId)

	err := RefundOrder(ctx, mockStorage, 554, 124)
	assert.ErrorIs(t, err, e.ErrIncorrectUserId)
}

func TestRefundOrder_SuccessfulRefund(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockStorage := mocks.NewFacadeMock(ctrl)

	mockStorage.RefundOrderMock.Expect(ctx, uint(88), uint(123)).Return(nil)

	err := RefundOrder(ctx, mockStorage, 88, 123)
	assert.NoError(t, err)
}
