package storage

import (
	"encoding/json"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"io"
	"os"
)

type ReturnStorageInterface interface {
	AddReturnToStorage(r *models.Return) error
	DeleteReturnFromStorage(id uint)
	IsConsist(id uint) bool
	GetReturn(id uint) (*models.Return, bool)
	GetReturnIDs() []uint
	ReadFromJSON() error
	WriteToJSON() error
}

type ReturnStorage struct {
	Data map[uint]*models.Return
	path string
}

func (r *ReturnStorage) Create() error {
	_, err := os.Create(r.path)
	return err
}

func NewReturnStorage() *ReturnStorage {
	return &ReturnStorage{Data: make(map[uint]*models.Return), path: "api/returns.json"}
}

func (rs *ReturnStorage) AddReturnToStorage(r *models.Return) error {
	_, ok := rs.Data[r.ID]
	if ok {
		return ErrAlrReturn
	} else {
		rs.Data[r.ID] = r
	}

	return nil
}

func (rs *ReturnStorage) DeleteReturnFromStorage(id uint) {
	delete(rs.Data, id)
}

// проверка на наличие
func (rs *ReturnStorage) IsConsist(id uint) bool {
	_, ok := rs.Data[id]
	return ok
}

// считываем с JSON-a
func (rs *ReturnStorage) ReadFromJSON() error {
	file, err := os.OpenFile(rs.path, os.O_RDONLY, 0666)
	if err != nil {
		return ErrOpenFile
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return ErrReadFile
	}

	if len(data) == 0 {
		rs.Data = make(map[uint]*models.Return)
		return nil
	}

	var i struct {
		Data map[uint]models.Return `json:"Data"`
	}

	err = json.Unmarshal(data, &i)
	if err != nil {
		return err
	}

	rs.Data = make(map[uint]*models.Return)
	for returnid, r := range i.Data {
		returnCopy := r
		rs.Data[returnid] = &returnCopy
	}

	return nil
}

func (rs *ReturnStorage) WriteToJSON() error {
	file, err := os.OpenFile("api/returns.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return ErrOpenFile
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent(" ", "  ")
	if err := encoder.Encode(rs); err != nil {
		return ErrEnocde
	}
	return nil
}

func (rs *ReturnStorage) GetReturn(id uint) (*models.Return, bool) {
	r, ok := rs.Data[id]
	return r, ok
}

func (rs *ReturnStorage) GetReturnIDs() []uint {
	var ids []uint
	for id := range rs.Data {
		ids = append(ids, id)
	}
	return ids
}
