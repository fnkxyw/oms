package orders

import (
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/mocks"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
	"testing"
	"time"
)

// этот тест я сделал немного по-другому паттерну, т.к c setupMock была проблема c Order.AcceptTime из-за вызова функций разница была около 1/100 миллисекунды
func TestAcceptOrder(t *testing.T) {
	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewOrderStorageInterfaceMock(ctrl)

	tests := []struct {
		name string
		args struct {
			s  storage.OrderStorageInterface
			or *models.Order
		}
		isOrderExist       bool // существует ли заказ
		expectAddToStorage bool // нужно ли добавлять заказ
		wantErr            assert.ErrorAssertionFunc
	}{
		{
			name: "successful acceptance",
			args: struct {
				s  storage.OrderStorageInterface
				or *models.Order
			}{
				s: mockStorage,
				or: &models.Order{
					ID:            1,
					KeepUntilDate: time.Now().Add(24 * time.Hour),
					State:         models.AcceptState,
				},
			},
			isOrderExist:       false,
			expectAddToStorage: true,
			wantErr:            assert.NoError,
		},
		{
			name: "order already exists",
			args: struct {
				s  storage.OrderStorageInterface
				or *models.Order
			}{
				s: mockStorage,
				or: &models.Order{
					ID:            1,
					KeepUntilDate: time.Now().Add(24 * time.Hour),
					State:         models.AcceptState,
				},
			},
			isOrderExist:       true,
			expectAddToStorage: false,
			wantErr:            assert.Error,
		},
		{
			name: "order has expired",
			args: struct {
				s  storage.OrderStorageInterface
				or *models.Order
			}{
				s: mockStorage,
				or: &models.Order{
					ID:            2,
					KeepUntilDate: time.Now().Add(-24 * time.Hour),
					State:         models.AcceptState,
				},
			},
			isOrderExist:       false,
			expectAddToStorage: false,
			wantErr:            assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockStorage.IsConsistMock.Expect(tt.args.or.ID).Return(tt.isOrderExist)

			if !tt.isOrderExist && tt.expectAddToStorage {
				mockStorage.AddOrderToStorageMock.Expect(tt.args.or)
			}

			err := AcceptOrder(tt.args.s, tt.args.or)
			tt.wantErr(t, err)
		})
	}
}

func TestPlaceOrder(t *testing.T) {
	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewOrderStorageInterfaceMock(ctrl)

	tests := []struct {
		name      string
		ids       []uint
		setupMock func()
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "empty ids list",
			ids:  []uint{},
			setupMock: func() {

			},
			wantErr: assert.Error,
		},
		{
			name: "no consist",
			ids:  []uint{1},
			setupMock: func() {
				mockStorage.GetOrderMock.Expect(uint(1)).Return(&models.Order{
					ID:            1,
					State:         models.NewState,
					KeepUntilDate: time.Now().Add(24 * time.Hour),
				}, false)
			},
			wantErr: assert.Error,
		},
		{
			name: "order already placed",
			ids:  []uint{1},
			setupMock: func() {
				//заказ уже размещен
				mockStorage.GetOrderMock.Expect(uint(1)).Return(&models.Order{
					ID:            1,
					State:         models.PlaceState,
					KeepUntilDate: time.Now().Add(24 * time.Hour),
				}, true)
			},
			wantErr: assert.Error,
		},
		{
			name: "order was deleted",
			ids:  []uint{2},
			setupMock: func() {
				// заказ был удалён
				mockStorage.GetOrderMock.Expect(uint(2)).Return(&models.Order{
					ID:            2,
					State:         models.SoftDelete,
					KeepUntilDate: time.Now().Add(24 * time.Hour),
				}, true)
			},
			wantErr: assert.Error,
		},
		{
			name: "order expired",
			ids:  []uint{3},
			setupMock: func() {
				mockStorage.GetOrderMock.Expect(uint(3)).Return(&models.Order{
					ID:            3,
					State:         models.AcceptState,
					KeepUntilDate: time.Now().Add(-24 * time.Hour),
				}, true)
			},
			wantErr: assert.Error,
		},
		{
			name: "successful order placement",
			ids:  []uint{4},
			setupMock: func() {
				// Успешный сценарий
				mockStorage.GetOrderMock.Expect(uint(4)).Return(&models.Order{
					ID:            4,
					State:         models.AcceptState,
					KeepUntilDate: time.Now().Add(24 * time.Hour),
				}, true)
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := PlaceOrder(mockStorage, tt.ids)

			tt.wantErr(t, err)
		})
	}
}

func TestListOrders(t *testing.T) {
	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewOrderStorageInterfaceMock(ctrl)

	tests := []struct {
		name      string
		id        uint
		n         int
		inPuP     bool
		setupMock func()
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name:  "no_orders_available",
			id:    1,
			n:     1,
			inPuP: false,
			setupMock: func() {
				mockStorage.IsConsistMock.Expect(uint(1)).Return(true)
				mockStorage.GetOrderIDsMock.Expect().Return([]uint{})
			},
			wantErr: assert.NoError,
		},
		{
			name:  "no consists",
			id:    1,
			n:     1,
			inPuP: false,
			setupMock: func() {
				mockStorage.IsConsistMock.Expect(uint(1)).Return(false)
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := ListOrders(mockStorage, tt.id, tt.n, tt.inPuP)
			tt.wantErr(t, err)

		})
	}
}

func TestReturnOrder(t *testing.T) {
	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewOrderStorageInterfaceMock(ctrl)

	type args struct {
		s  storage.OrderStorageInterface
		id uint
	}
	tests := []struct {
		name      string
		args      args
		setupMock func()
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "successful return returned state",
			args: args{
				s:  mockStorage,
				id: 1,
			},
			setupMock: func() {
				mockStorage.GetOrderMock.Expect(uint(1)).Return(&models.Order{
					ID:     1,
					UserID: 1,
					State:  models.ReturnedState,
				}, true)
			},
			wantErr: assert.NoError,
		},
		{
			name: "successful return date and accept state",
			args: args{
				s:  mockStorage,
				id: 1,
			},
			setupMock: func() {
				mockStorage.GetOrderMock.Expect(uint(1)).Return(&models.Order{
					ID:            1,
					UserID:        1,
					State:         models.AcceptState,
					KeepUntilDate: time.Now().Add(-24 * time.Hour),
				}, true)
			},
			wantErr: assert.NoError,
		},
		{
			name: "no order found",
			args: args{
				s:  mockStorage,
				id: 1,
			},
			setupMock: func() {
				mockStorage.GetOrderMock.Expect(uint(1)).Return(&models.Order{}, false)
			},
			wantErr: assert.Error,
		},
		{
			name: "order can't be returned",
			args: args{
				s:  mockStorage,
				id: 1,
			},
			setupMock: func() {
				mockStorage.GetOrderMock.Expect(uint(1)).Return(&models.Order{
					ID:            1,
					UserID:        1,
					State:         models.AcceptState,
					KeepUntilDate: time.Now().Add(24 * time.Hour),
				}, true)
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := ReturnOrder(tt.args.s, tt.args.id)

			tt.wantErr(t, err)
		})
	}
}
