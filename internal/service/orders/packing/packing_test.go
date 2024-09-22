package packing_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders/packing"
)

func TestPackaging(t *testing.T) {
	tests := []struct {
		name string
		args struct {
			order        *models.Order
			packType     string
			needWrapping bool
		}
		wantErr       assert.ErrorAssertionFunc
		expectedPrice int
	}{
		{
			name: "box packaging with wrap",
			args: struct {
				order        *models.Order
				packType     string
				needWrapping bool
			}{
				order:        &models.Order{Weight: 20, Price: 100},
				packType:     "box",
				needWrapping: true,
			},
			wantErr:       assert.NoError,
			expectedPrice: 121,
		},
		{
			name: "box packaging without wrap",
			args: struct {
				order        *models.Order
				packType     string
				needWrapping bool
			}{
				order:        &models.Order{Weight: 20, Price: 100},
				packType:     "box",
				needWrapping: false,
			},
			wantErr:       assert.NoError,
			expectedPrice: 120,
		},
		{
			name: "bundle packaging with wrap",
			args: struct {
				order        *models.Order
				packType     string
				needWrapping bool
			}{
				order:        &models.Order{Weight: 5, Price: 100},
				packType:     "bundle",
				needWrapping: true,
			},
			wantErr:       assert.NoError,
			expectedPrice: 106,
		},
		{
			name: "invalid weight for box",
			args: struct {
				order        *models.Order
				packType     string
				needWrapping bool
			}{
				order:        &models.Order{Weight: 31, Price: 100},
				packType:     "box",
				needWrapping: false,
			},
			wantErr:       assert.Error,
			expectedPrice: 100,
		},
		{
			name: "invalid package type",
			args: struct {
				order        *models.Order
				packType     string
				needWrapping bool
			}{
				order:        &models.Order{Weight: 5, Price: 100},
				packType:     "invalid",
				needWrapping: false,
			},
			wantErr:       assert.Error,
			expectedPrice: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := packing.Packing(tt.args.order, tt.args.packType, tt.args.needWrapping)

			tt.wantErr(t, err)

			assert.Equal(t, tt.expectedPrice, tt.args.order.Price, "expected price to be %v, got %v", tt.expectedPrice, tt.args.order.Price)
		})
	}
}
