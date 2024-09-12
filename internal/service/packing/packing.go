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
	if o.IsPackaged == true {
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
	if o.IsPackaged == true {
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

// процесс упаковки
func Packing(o *models.Order, pack string) error {
	var packager Packager

	switch pack {
	case "box":
		packager = &BoxPackaging{}
	case "bundle":
		packager = &BundlePackaging{}
	case "wrap":
		packager = &WrapPackaging{}
	default:
		return ErrorInvalidType
	}

	return packager.Pack(o)
}
