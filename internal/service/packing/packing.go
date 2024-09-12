package packing

import (
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/service/packing/controller"
)

// интерфейс упаковки
type Packager interface {
	Pack(o *models.Order) error
}

type BoxPackaging struct {
}

func (b *BoxPackaging) Pack(o *models.Order) error {
	if o.Weight > 30 {
		return ErrorWeightBox
	}
	if o.IsPackaged {
		return ErrorIsPackaged
	}
	o.Price += 20
	o.IsPackaged = true
	ans, err := controller.AddWrap()
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
		return ErrorWeightBundle
	}
	if o.IsPackaged {
		return ErrorIsPackaged
	}

	o.IsPackaged = true
	o.Price += 5

	ans, err := controller.AddWrap()
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

func GetPackager(pack string) (Packager, error) {
	switch pack {
	case "box":
		return &BoxPackaging{}, nil
	case "bundle":
		return &BundlePackaging{}, nil
	case "wrap":
		return &WrapPackaging{}, nil
	default:
		return nil, ErrorInvalidType
	}
}

func Packing(o *models.Order, pack string) error {
	packager, err := GetPackager(pack)
	if err != nil {
		return err
	}
	return packager.Pack(o)
}
