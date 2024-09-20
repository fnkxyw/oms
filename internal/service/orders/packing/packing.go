package packing

import (
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders/packing/controller"
)

// интерфейс упаковки
type Packager interface {
	Pack(o *models.Order) error
}

type BoxPackaging struct {
}

func (b *BoxPackaging) Pack(o *models.Order) error {
	if o.Weight > 30 {
		return ErrWeightBox
	}

	o.Price += 20

	ans, err := controller.GetWrapAdder().AddWrap()
	if err != nil {
		return err
	}
	if ans {
		Packing(o, "wrap")
	}
	return nil
}

type BundlePackaging struct {
}

func (b *BundlePackaging) Pack(o *models.Order) error {
	if o.Weight > 10 {
		return ErrWeightBundle
	}

	o.Price += 5

	ans, err := controller.GetWrapAdder().AddWrap()
	if err != nil {
		return err
	}
	if ans {
		Packing(o, "wrap")
	}

	return nil
}

type WrapPackaging struct {
}

func (w *WrapPackaging) Pack(o *models.Order) error {
	o.Price += 1
	return nil
}

// рациональное решение в случае если в дальнейшем понадобиться добавить еще одну упаковку
func GetPackager(pack string) (Packager, error) {
	switch pack {
	case "box":
		return &BoxPackaging{}, nil
	case "bundle":
		return &BundlePackaging{}, nil
	case "wrap":
		return &WrapPackaging{}, nil
	default:
		return nil, ErrInvalidType
	}
}

// упаковка
func Packing(o *models.Order, pack string) error {
	packager, err := GetPackager(pack)
	if err != nil {
		return err
	}
	return packager.Pack(o)
}
