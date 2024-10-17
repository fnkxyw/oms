package packing

import (
	"errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
)

type PackageType int32

const (
	PackageType_PACKAGE_UNKNOWN PackageType = 0
	PackageType_BOX             PackageType = 1
	PackageType_BUNDLE          PackageType = 2
	PackageType_WRAP            PackageType = 3
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
		err := Packing(o, PackageType_WRAP, false)
		if err != nil {
			return errors.New("Packing in wrap error ")
		}
	}

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
		err := Packing(o, PackageType_WRAP, false)
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
func GetPackager(pack PackageType) (Packager, error) {
	switch pack {
	case PackageType_BOX:
		return &BoxPackaging{}, nil
	case PackageType_BUNDLE:
		return &BundlePackaging{}, nil
	case PackageType_WRAP:
		return &WrapPackaging{}, nil
	default:
		return nil, ErrInvalidType
	}
}

// упаковка
func Packing(o *models.Order, pack PackageType, needWrapping bool) error {
	packager, err := GetPackager(pack)
	if err != nil {
		return err
	}
	return packager.Pack(o, needWrapping)
}
