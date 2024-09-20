package orders

import (
	"github.com/gojuno/minimock/v3"
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
		wantErr            bool
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
			wantErr:            false,
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
			wantErr:            true,
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
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockStorage.IsConsistMock.Expect(tt.args.or.ID).Return(tt.isOrderExist)

			if !tt.isOrderExist && tt.expectAddToStorage {
				mockStorage.AddOrderToStorageMock.Expect(tt.args.or)
			}

			err := AcceptOrder(tt.args.s, tt.args.or)

			if (err != nil) != tt.wantErr {
				t.Errorf("AcceptOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantErr   bool
	}{
		{
			name: "empty ids list",
			ids:  []uint{},
			setupMock: func() {

			},
			wantErr: true,
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
			wantErr: true,
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
			wantErr: true,
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
			wantErr: true,
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
			wantErr: true,
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
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := PlaceOrder(mockStorage, tt.ids)

			if (err != nil) != tt.wantErr {
				t.Errorf("PlaceOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantErr   bool
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
			wantErr: false,
		},
		{
			name:  "no consists",
			id:    1,
			n:     1,
			inPuP: false,
			setupMock: func() {
				mockStorage.IsConsistMock.Expect(uint(1)).Return(false)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := ListOrders(mockStorage, tt.id, tt.n, tt.inPuP)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListOrders() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestRefundOrder(t *testing.T) {
	ctrl := minimock.NewController(t)
	defer ctrl.Finish()
	mockOrderStorage := mocks.NewOrderStorageInterfaceMock(ctrl)
	mockReturnStorage := mocks.NewReturnStorageInterfaceMock(ctrl)

	userId := uint(123)

	tests := []struct {
		name      string
		orderId   uint
		userId    uint
		setupMock func()
		wantErr   bool
	}{
		{
			name:    "order does not exist",
			orderId: 1,
			userId:  userId,
			setupMock: func() {
				mockOrderStorage.GetOrderMock.Expect(uint(1)).Return(nil, false)
			},
			wantErr: true,
		},
		{
			name:    "order is not in PlaceState",
			orderId: 2,
			userId:  userId,
			setupMock: func() {
				mockOrderStorage.GetOrderMock.Expect(uint(2)).Return(&models.Order{
					ID:        2,
					State:     models.AcceptState,
					UserID:    userId,
					PlaceDate: time.Now(),
				}, true)
			},
			wantErr: true,
		},
		{
			name:    "order refund time expired",
			orderId: 3,
			userId:  userId,
			setupMock: func() {
				mockOrderStorage.GetOrderMock.Expect(uint(3)).Return(&models.Order{
					ID:        3,
					State:     models.PlaceState,
					UserID:    userId,
					PlaceDate: time.Now().AddDate(0, 0, -3),
				}, true)
			},
			wantErr: true,
		},
		{
			name:    "incorrect user ID",
			orderId: 4,
			userId:  userId + 1,
			setupMock: func() {
				mockOrderStorage.GetOrderMock.Expect(uint(4)).Return(&models.Order{
					ID:        4,
					State:     models.PlaceState,
					UserID:    userId,
					PlaceDate: time.Now(),
				}, true)
			},
			wantErr: true,
		},
		{
			name:    "successful refund",
			orderId: 5,
			userId:  userId,
			setupMock: func() {
				mockOrderStorage.GetOrderMock.Expect(uint(5)).Return(&models.Order{
					ID:        5,
					State:     models.PlaceState,
					UserID:    userId,
					PlaceDate: time.Now(),
				}, true)
				mockReturnStorage.AddReturnToStorageMock.Expect(&models.Return{
					ID:     5,
					UserID: userId,
				}).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := RefundOrder(mockReturnStorage, mockOrderStorage, tt.orderId, tt.userId)

			if (err != nil) != tt.wantErr {
				t.Errorf("RefundOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
