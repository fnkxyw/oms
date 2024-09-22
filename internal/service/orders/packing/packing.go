package packing

import (
	"errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
)

// интерфейс упаковки
type Packager interface {
	Pack(o *models.Order, needWrapping bool) error
}

type BoxPackaging struct {
}

func (b *BoxPackaging) Pack(o *models.Order, needWrapping bool) error {
	if o.Weight > 30 {
		return ErrWeightBox
	}

	o.Price += 20

	if needWrapping {
		err := Packing(o, "wrap", false)
		if err != nil {
			return errors.New("Packing in wrap error ")
		}
	}
	//ans, err := controller.GetWrapAdder().AddWrap()
	//if err != nil {
	//	return err
	//}
	//if ans {
	//	err = Packing(o, "wrap")
	//	if err != nil {
	//		return err
	//	}
	//}
	return nil
}

type BundlePackaging struct {
}

func (b *BundlePackaging) Pack(o *models.Order, needWrapping bool) error {
	if o.Weight > 10 {
		return ErrWeightBundle
	}

	o.Price += 5

	if needWrapping {
		err := Packing(o, "wrap", false)
		if err != nil {
			return errors.New("Packing in wrap error ")
		}
	}

	return nil
}

type WrapPackaging struct {
}

func (w *WrapPackaging) Pack(o *models.Order, needWrapping bool) error {
	if needWrapping {
		return errors.New("you can`t wrapping wrap")
	}
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
func Packing(o *models.Order, pack string, needWrapping bool) error {
	packager, err := GetPackager(pack)
	if err != nil {
		return err
	}
	return packager.Pack(o, needWrapping)
}
