package returns

import (
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/mocks"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
	"testing"
	"time"
)

// тут я сделал так, что в начале теста определил поведение мока, потому что функция легко покрывается такими тестовыми данными
func TestListReturns(t *testing.T) {
	ctrl := minimock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewReturnStorageInterfaceMock(ctrl)

	// Структура для аргументов
	type args struct {
		s     storage.ReturnStorageInterface
		limit int
		page  int
	}

	tests := []struct {
		name      string
		args      args
		wantErr   assert.ErrorAssertionFunc
		setupMock func()
	}{
		{
			name: "valid_returns_with_pagination",
			args: args{
				s:     mockStorage,
				limit: 2,
				page:  2,
			},
			setupMock: func() {
				mockStorage.GetReturnIDsMock.Expect().Return([]uint{1, 2, 3, 4, 5})
				mockStorage.GetReturnMock.When(uint(1)).Then(&models.Return{ID: 1, UserID: 1001}, true)
				mockStorage.GetReturnMock.When(uint(2)).Then(&models.Return{ID: 2, UserID: 1002}, true)
				mockStorage.GetReturnMock.When(uint(3)).Then(&models.Return{ID: 3, UserID: 1003}, true)
				mockStorage.GetReturnMock.When(uint(4)).Then(&models.Return{ID: 4, UserID: 1004}, true)
				mockStorage.GetReturnMock.When(uint(5)).Then(&models.Return{ID: 5, UserID: 1005}, true)
			},
			wantErr: assert.NoError,
		},
		{
			name: "invalid_page_number",
			args: args{
				s:     mockStorage,
				limit: 2,
				page:  -1,
			},
			setupMock: func() {
				mockStorage.GetReturnIDsMock.Expect().Return([]uint{1, 2, 3, 4, 5})
				mockStorage.GetReturnMock.When(uint(1)).Then(&models.Return{ID: 1, UserID: 1001}, true)
				mockStorage.GetReturnMock.When(uint(2)).Then(&models.Return{ID: 2, UserID: 1002}, true)
				mockStorage.GetReturnMock.When(uint(3)).Then(&models.Return{ID: 3, UserID: 1003}, true)
				mockStorage.GetReturnMock.When(uint(4)).Then(&models.Return{ID: 4, UserID: 1004}, true)
				mockStorage.GetReturnMock.When(uint(5)).Then(&models.Return{ID: 5, UserID: 1005}, true)
			},
			wantErr: assert.Error,
		},

		{
			name: "no_more_items",
			args: args{
				s:     mockStorage,
				limit: 2,
				page:  5, // запрашиваем несуществующую страницу
			},
			setupMock: func() {
				mockStorage.GetReturnIDsMock.Expect().Return([]uint{1, 2, 3, 4, 5})
				mockStorage.GetReturnMock.When(uint(1)).Then(&models.Return{ID: 1, UserID: 1001}, true)
				mockStorage.GetReturnMock.When(uint(2)).Then(&models.Return{ID: 2, UserID: 1002}, true)
				mockStorage.GetReturnMock.When(uint(3)).Then(&models.Return{ID: 3, UserID: 1003}, true)
				mockStorage.GetReturnMock.When(uint(4)).Then(&models.Return{ID: 4, UserID: 1004}, true)
				mockStorage.GetReturnMock.When(uint(5)).Then(&models.Return{ID: 5, UserID: 1005}, true)
			},
			wantErr: assert.Error,
		},
		{
			name: "invalid_limit_number",
			args: args{
				s:     mockStorage,
				limit: 0, // неправильный лимит
				page:  1,
			},
			setupMock: func() {
				mockStorage.GetReturnIDsMock.Expect().Return([]uint{1, 2, 3, 4, 5})
				mockStorage.GetReturnMock.When(uint(1)).Then(&models.Return{ID: 1, UserID: 1001}, true)
				mockStorage.GetReturnMock.When(uint(2)).Then(&models.Return{ID: 2, UserID: 1002}, true)
				mockStorage.GetReturnMock.When(uint(3)).Then(&models.Return{ID: 3, UserID: 1003}, true)
				mockStorage.GetReturnMock.When(uint(4)).Then(&models.Return{ID: 4, UserID: 1004}, true)
				mockStorage.GetReturnMock.When(uint(5)).Then(&models.Return{ID: 5, UserID: 1005}, true)
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setupMock()
			err := ListReturns(mockStorage, tt.args.limit, tt.args.page)
			tt.wantErr(t, err)
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
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name:    "order does not exist",
			orderId: 1,
			userId:  userId,
			setupMock: func() {
				mockOrderStorage.GetOrderMock.Expect(uint(1)).Return(nil, false)
			},
			wantErr: assert.Error,
		},
		{
			name:    "order is not in PlaceState",
			orderId: 2,
			userId:  userId,
			setupMock: func() {
				mockOrderStorage.GetOrderMock.When(uint(2)).Then(&models.Order{
					ID:        2,
					State:     models.AcceptState,
					UserID:    userId,
					PlaceDate: time.Now(),
				}, true)
			},
			wantErr: assert.Error,
		},
		{
			name:    "order refund time expired",
			orderId: 3,
			userId:  userId,
			setupMock: func() {
				mockOrderStorage.GetOrderMock.When(uint(3)).Then(&models.Order{
					ID:        3,
					State:     models.PlaceState,
					UserID:    userId,
					PlaceDate: time.Now().AddDate(0, 0, -3),
				}, true)
			},
			wantErr: assert.Error,
		},
		{
			name:    "incorrect user ID",
			orderId: 4,
			userId:  userId + 1,
			setupMock: func() {
				mockOrderStorage.GetOrderMock.When(uint(4)).Then(&models.Order{
					ID:        4,
					State:     models.PlaceState,
					UserID:    userId,
					PlaceDate: time.Now(),
				}, true)
			},
			wantErr: assert.Error,
		},
		{
			name:    "successful refund",
			orderId: 5,
			userId:  userId,
			setupMock: func() {
				mockOrderStorage.GetOrderMock.When(uint(5)).Then(&models.Order{
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
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setupMock()

			err := RefundOrder(mockReturnStorage, mockOrderStorage, tt.orderId, tt.userId)
			tt.wantErr(t, err)
		})
	}
}
