package packing_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders/packing"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders/packing/controller"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders/packing/mocks"
)

func TestPackaging(t *testing.T) {
	tests := []struct {
		name string
		args struct {
			order    *models.Order
			packType string
		}
		setupMock     func()
		wantErr       assert.ErrorAssertionFunc
		expectedPrice int
	}{
		{
			name: "box packaging with wrap",
			args: struct {
				order    *models.Order
				packType string
			}{
				order:    &models.Order{Weight: 20, Price: 100},
				packType: "box",
			},
			setupMock: func() {
				mockWrapAdder := mocks.NewWrapAdderMock(t)
				mockWrapAdder.AddWrapMock.Expect().Return(true, nil)
				controller.SetWrapAdder(mockWrapAdder)
			},
			wantErr:       assert.NoError,
			expectedPrice: 121,
		},
		{
			name: "Box Packaging Without Wrap",
			args: struct {
				order    *models.Order
				packType string
			}{
				order:    &models.Order{Weight: 20, Price: 100},
				packType: "box",
			},
			setupMock: func() {
				mockWrapAdder := mocks.NewWrapAdderMock(t)
				mockWrapAdder.AddWrapMock.Expect().Return(false, nil)
				controller.SetWrapAdder(mockWrapAdder)
			},
			wantErr:       assert.NoError,
			expectedPrice: 120,
		},
		{
			name: "Bundle Packaging With Wrap",
			args: struct {
				order    *models.Order
				packType string
			}{
				order:    &models.Order{Weight: 5, Price: 100},
				packType: "bundle",
			},
			setupMock: func() {
				mockWrapAdder := mocks.NewWrapAdderMock(t)
				mockWrapAdder.AddWrapMock.Expect().Return(true, nil)
				controller.SetWrapAdder(mockWrapAdder)
			},
			wantErr:       assert.NoError,
			expectedPrice: 106,
		},
		{
			name: "Invalid Weight Box",
			args: struct {
				order    *models.Order
				packType string
			}{
				order:    &models.Order{Weight: 31, Price: 100},
				packType: "box",
			},
			setupMock:     func() {},
			wantErr:       assert.Error,
			expectedPrice: 100,
		},
		{
			name: "Invalid Package Type",
			args: struct {
				order    *models.Order
				packType string
			}{
				order:    &models.Order{Weight: 5, Price: 100},
				packType: "invalid",
			},
			setupMock:     func() {},
			wantErr:       assert.Error,
			expectedPrice: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := packing.Packing(tt.args.order, tt.args.packType)
			tt.wantErr(t, err)
			if err != nil && tt.args.order.Price != tt.expectedPrice {
				t.Fatalf("expected price to be %v, got %v", tt.expectedPrice, tt.args.order.Price)
			}
		})
	}
}
